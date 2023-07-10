package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepo struct {
	db *postgres.Postgres
}

func NewOrderRepo(db *postgres.Postgres) *OrderRepo {
	return &OrderRepo{db: db}
}

func (repo *OrderRepo) AddOrder(ctx context.Context, order *entities.Order) error {
	query := `
		INSERT INTO orders(user_id, order_num, status, uploaded)
		VALUES($1, $2, $3, now())
	`

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, order.UserID, order.Number, order.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code != pgerrcode.UniqueViolation {
			return err
		}

		existOrder, err := repo.OrderByNumber(ctx, order.Number)
		if err != nil {
			return err
		}

		if order.UserID == existOrder.UserID {
			return entities.ErrOrderAlreadyAdded
		}

		return entities.ErrOrderAddedByOther
	}

	return tx.Commit()
}

func (repo *OrderRepo) OrderByNumber(ctx context.Context, number string) (*entities.Order, error) {
	query := `
		SELECT user_id, order_num, status, accrual, uploaded
		FROM orders 
		WHERE order_num = $1
	`

	var (
		accrual sql.NullInt32
		order   = &entities.Order{}
	)

	err := repo.db.QueryRowContext(ctx, query, number).Scan(
		&order.UserID, &order.Number, &order.Status, &accrual, &order.UploadedAt,
	)
	if err != nil {
		return nil, err
	}

	order.Accrual = int(accrual.Int32)

	return order, nil
}

func (repo *OrderRepo) Orders(ctx context.Context, userID int64) ([]entities.Order, error) {
	query := `
		SELECT order_num, status, accrual, uploaded
		FROM orders 
		WHERE user_id = $1
	`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]entities.Order, 0)

	for rows.Next() {
		var (
			accrual sql.NullInt32
			order   = entities.Order{UserID: userID}
		)

		err = rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		order.Accrual = int(accrual.Int32)

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
