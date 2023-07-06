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

type UserRepo struct {
	db *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) AddUser(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users(login, password, salt) VALUES($1, $2, $3)
		RETURNING id
	`

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var id int64

	err = tx.QueryRowContext(ctx, query, user.Login, user.EncryptedPassword, user.Salt).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return entities.ErrUserAlreadyExists
		}

		return err
	}

	user.ID = id

	return tx.Commit()
}

func (repo *UserRepo) User(ctx context.Context, user *entities.User) error {
	query := `
		SELECT 
			id, password, salt 
		FROM users
		WHERE login = $1
	`

	err := repo.db.QueryRowContext(ctx, query, user.Login).Scan(&user.ID, &user.EncryptedPassword, &user.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.ErrInvalidLoginPassword
		}

		return err
	}

	return nil
}
