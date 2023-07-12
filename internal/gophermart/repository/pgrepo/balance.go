package pgrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/postgres"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type BalanceRepo struct {
	db *postgres.Postgres
}

func NewBalanceRepo(db *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{db: db}
}

func (repo *BalanceRepo) Balance(ctx context.Context, userID int64) (entities.Balance, error) {
	query := `
		SELECT ub.balance, COALESCE(ubl.withdrawals, 0)
		FROM user_balance ub
		LEFT JOIN (
			SELECT user_id, operation, sum(sum) AS withdrawals
			FROM user_balance_log
			WHERE operation = 'withdrawal'
			GROUP BY user_id, operation
		) ubl ON ub.user_id = ubl.user_id
		WHERE ub.user_id = $1
	`

	balance := entities.Balance{UserID: userID}

	err := repo.db.QueryRowContext(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return entities.Balance{}, err
	}

	return balance, nil
}

func (repo *BalanceRepo) ChangeBalance(ctx context.Context, change *entities.BalanceChange) error {
	query1 := `
		UPDATE user_balance
		SET balance = balance %s $1
		WHERE user_id = $2
	`

	if change.Operation == entities.BalanceOperationWithdrawal {
		query1 = fmt.Sprintf(query1, "-")
	} else {
		query1 = fmt.Sprintf(query1, "+")
	}

	query2 := `
		INSERT INTO user_balance_log(user_id, processed, operation, order_num, sum)
		VALUES ($1, now(), $2, $3, $4)
	`

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query1, change.Sum, change.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.CheckViolation {
			return entities.ErrNotEnoughFunds
		}
	}

	_, err = tx.ExecContext(ctx, query2, change.UserID, change.Operation, change.Order, change.Sum)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *BalanceRepo) Withdrawals(ctx context.Context, userID int64) ([]entities.BalanceChange, error) {
	query := `
		SELECT order_num, sum, processed
		FROM user_balance_log
		WHERE user_id = $1
		ORDER BY processed ASC
	`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	withdrawals := make([]entities.BalanceChange, 0)

	for rows.Next() {
		withdrawal := entities.BalanceChange{
			UserID:    userID,
			Operation: entities.BalanceOperationWithdrawal,
		}

		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
