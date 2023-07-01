package pgrepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

type PgRepo struct {
	db *sql.DB
}

func NewPgRepo(dsn, migrations string) (*PgRepo, error) {
	dbPool, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping()
	if err != nil {
		return nil, err
	}

	repo := &PgRepo{
		db: dbPool,
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

func (repo *PgRepo) Register() error {
	return nil
}

func (repo *PgRepo) Close() error {
	return repo.db.Close()
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
