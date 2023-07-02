package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
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

func (repo *PgRepo) CreateUser() error {
	return nil
}

func (repo *PgRepo) Login() error {
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
