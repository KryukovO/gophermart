package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type BalanceController struct {
	balance usecases.Balance
	mw      *middleware.Manager
	logger  *log.Logger
}

func NewBalanceController(
	balance usecases.Balance, mwManager *middleware.Manager, logger *log.Logger,
) (*BalanceController, error) {
	if balance == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &BalanceController{
		balance: balance,
		mw:      mwManager,
		logger:  controllerLogger,
	}, nil
}

func (c *BalanceController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	router.Add(http.MethodGet, "/api/user/balance", c.mw.AuthenticationMiddleware(c.balanceHandler))
	router.Add(http.MethodPost, "/api/user/balance/withdraw", c.mw.AuthenticationMiddleware(c.withdrawHandler))
	router.Add(http.MethodGet, "/api/user/withdrawals", c.mw.AuthenticationMiddleware(c.withdrawalsHandler))

	return nil
}

func (c *BalanceController) balanceHandler(e echo.Context) error {
	uuid := e.Get("uuid")
	if uuid == nil {
		uuid = ""
	}

	userID := e.Get("userID")

	user, ok := userID.(int64)
	if !ok {
		return e.NoContent(http.StatusUnauthorized)
	}

	balance, err := c.balance.Balance(e.Request().Context(), user)
	if err != nil {
		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	return e.JSON(http.StatusOK, &balance)
}

func (c *BalanceController) withdrawHandler(e echo.Context) error {
	uuid := e.Get("uuid")
	if uuid == nil {
		uuid = ""
	}

	userID := e.Get("userID")

	user, ok := userID.(int64)
	if !ok {
		return e.NoContent(http.StatusUnauthorized)
	}

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	var change entities.BalanceChange

	err = json.Unmarshal(body, &change)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	c.logger.Debugf("Request body: %+v", change)

	change.UserID = user
	change.Operation = entities.BalanceOperationWithdrawal

	err = c.balance.ChangeBalance(e.Request().Context(), &change)
	if err != nil {
		if errors.Is(err, entities.ErrNotEnoughFunds) {
			return e.NoContent(http.StatusPaymentRequired)
		}

		if errors.Is(err, entities.ErrInvalidOrderNumber) {
			return e.NoContent(http.StatusUnprocessableEntity)
		}

		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	return e.NoContent(http.StatusOK)
}

func (c *BalanceController) withdrawalsHandler(e echo.Context) error {
	uuid := e.Get("uuid")
	if uuid == nil {
		uuid = ""
	}

	userID := e.Get("userID")

	user, ok := userID.(int64)
	if !ok {
		return e.NoContent(http.StatusUnauthorized)
	}

	withdrawals, err := c.balance.Withdrawals(e.Request().Context(), user)
	if err != nil {
		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	if len(withdrawals) == 0 {
		return e.NoContent(http.StatusNoContent)
	}

	return e.JSON(http.StatusOK, withdrawals)
}
