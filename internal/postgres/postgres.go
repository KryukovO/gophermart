package postgres

import (
	"context"
	"database/sql"
)

type Postgres struct {
	*sql.DB
}

func NewPostgres(ctx context.Context, dsn, migrations string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()

		return nil, err
	}

	err = RunMigrations(dsn, migrations)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: db,
	}, nil
}
