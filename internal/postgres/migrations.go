package postgres

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigrations(dsn, migrations string) error {
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
