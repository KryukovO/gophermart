package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/KryukovO/gophermart/internal/entities"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgRepo struct {
	db *sql.DB
}

func NewPgRepo(ctx context.Context, dsn, migrations string) (*PgRepo, error) {
	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	repo := &PgRepo{
		db: database,
	}

	err = repo.Ping(ctx)
	if err != nil {
		return nil, err
	}

	err = repo.runMigrations(dsn, migrations)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *PgRepo) runMigrations(dsn, migrations string) error {
	migration, err := migrate.New(
		fmt.Sprintf("file://%s", migrations),
		dsn,
	)
	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (repo *PgRepo) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

func (repo *PgRepo) Close() error {
	return repo.db.Close()
}

func (repo *PgRepo) Register(ctx context.Context, user *entities.User) error {
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

func (repo *PgRepo) Login(ctx context.Context, user *entities.User) error {
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

func (repo *PgRepo) Orders() error {
	return nil
}

func (repo *PgRepo) AddOrder() error {
	return nil
}

func (repo *PgRepo) Balance() error {
	return nil
}

func (repo *PgRepo) Withdraw() error {
	return nil
}

func (repo *PgRepo) Withdrawals() error {
	return nil
}
