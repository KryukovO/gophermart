package main

import (
	"flag"

	"github.com/KryukovO/gophermart/internal/gophermart"
	"github.com/KryukovO/gophermart/internal/gophermart/config"

	"github.com/caarlos0/env"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

const (
	address        = ":8081"
	dsn            = ""
	accrualAddress = ""

	shutdownTimeout = 10
	migrations      = "sql/migrations"
)

func main() {
	cfg := new(config.Config)

	flag.StringVar(&cfg.Address, "a", address, "Address to run HTTP server")
	flag.StringVar(&cfg.DSN, "d", dsn, "URI to database")
	flag.StringVar(&cfg.AccrualAddress, "r", accrualAddress, "Accrual system address")
	flag.UintVar(&cfg.ShutdownTimeout, "shutdown", shutdownTimeout, "Server shutdown timeout")
	flag.StringVar(&cfg.Migrations, "migrations", migrations, "Directory of database migration files")
	flag.Parse()

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05 Z07:00",
	})

	err := env.Parse(cfg)
	if err != nil {
		logger.Fatalf("Env parsing error: %s. Exit(1)", err.Error())
	}

	if err := gophermart.Run(cfg, logger); err != nil {
		logger.Fatalf("Gophermart service error: %s. Exit(1)", err)
	}
}
