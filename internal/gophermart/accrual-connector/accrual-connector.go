package accrualconnector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedStatus = errors.New("unexpected response status")

type AccrualConnector struct {
	accrualAddr string
	workers     uint
	interval    time.Duration
	order       usecases.Order
	balance     usecases.Balance
	logger      *log.Logger
	close       chan struct{}
}

func NewAccrualConnector(
	accrualAddr string, workers uint, interval time.Duration,
	order usecases.Order, balance usecases.Balance,
	logger *log.Logger,
) *AccrualConnector {
	return &AccrualConnector{
		accrualAddr: accrualAddr,
		workers:     workers,
		interval:    interval,
		order:       order,
		balance:     balance,
		logger:      logger,
		close:       make(chan struct{}),
	}
}

func (connector *AccrualConnector) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-connector.close:
			return
		case <-time.After(connector.interval):
		}

		connector.logger.Info(ctx)

		orders, err := connector.order.ProcessableOrders(ctx)
		if err != nil {
			connector.logger.Errorf("AccrualConnector error: %s", err)
		}

		tasks := connector.generateOrderTasks(ctx, orders)

		group, gCtx := errgroup.WithContext(ctx)

		for w := 0; w < int(connector.workers); w++ {
			group.Go(func() error {
				return connector.orderTaskWorker(gCtx, tasks)
			})
		}

		err = group.Wait()
		if err != nil {
			connector.logger.Errorf("AccrualConnector error: %s", err)
		}
	}
}

func (connector *AccrualConnector) Shutdown(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case connector.close <- struct{}{}:
		return
	}
}

func (connector *AccrualConnector) generateOrderTasks(
	ctx context.Context, orders []entities.Order,
) chan entities.Order {
	outCh := make(chan entities.Order, connector.workers)

	go func() {
		defer close(outCh)

		for _, order := range orders {
			ord := order

			select {
			case <-ctx.Done():
				return
			case outCh <- ord:
			}
		}
	}()

	return outCh
}

func (connector *AccrualConnector) orderTaskWorker(ctx context.Context, tasks <-chan entities.Order) error {
	client := http.Client{}

	for order := range tasks {
		select {
		case <-ctx.Done():
			return nil
		default:
			accrualOrder, err := connector.doRequest(ctx, client, order.Number)
			if err != nil {
				return err
			}

			order.Status = entities.AccrualToOrderStatus(accrualOrder.Status)
			order.Accrual = accrualOrder.Accrual

			err = connector.order.UpdateOrder(ctx, &order)
			if err != nil {
				return err
			}

			if order.Status == entities.OrderStatusProcessed {
				balanceChange := entities.BalanceChange{
					UserID:    order.UserID,
					Operation: entities.BalanceOperationRefill,
					Order:     order.Number,
					Sum:       order.Accrual,
				}

				err := connector.balance.ChangeBalance(ctx, &balanceChange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (connector *AccrualConnector) doRequest(
	ctx context.Context, client http.Client, order string,
) (entities.AccrualOrder, error) {
	url := fmt.Sprintf("%s/api/orders/%s", connector.accrualAddr, order)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entities.AccrualOrder{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return entities.AccrualOrder{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entities.AccrualOrder{}, fmt.Errorf("%s: %w", resp.Status, ErrUnexpectedStatus)
	}

	var accrualOrder entities.AccrualOrder

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.AccrualOrder{}, err
	}

	err = json.Unmarshal(body, &accrualOrder)
	if err != nil {
		return entities.AccrualOrder{}, err
	}

	return accrualOrder, nil
}
