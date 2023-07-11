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

	secretKey         = ""
	repositoryTimeout = 3
	shutdownTimeout   = 10
	migrations        = "sql/migrations"
	userTokenTTL      = 30
	accrualWorkers    = 3
	accrualInterval   = 10
	accrualShutdown   = 3
)

func main() {
	cfg := new(config.Config)

	flag.StringVar(&cfg.Address, "a", address, "Address to run HTTP server")
	flag.StringVar(&cfg.DSN, "d", dsn, "URI to database")
	flag.StringVar(&cfg.AccrualAddress, "r", accrualAddress, "Accrual system address")

	flag.StringVar(&cfg.SecretKey, "secret", secretKey, "Authorization token encryption key")
	flag.UintVar(&cfg.RepositioryTimeout, "timeout", repositoryTimeout, "Repository connection timeout, sec")
	flag.UintVar(&cfg.ShutdownTimeout, "shutdown", shutdownTimeout, "Server shutdown timeout, sec")
	flag.StringVar(&cfg.Migrations, "migrations", migrations, "Directory of database migration files")
	flag.UintVar(&cfg.UserTokenTTL, "userttl", userTokenTTL, "User token lifetime, min")
	flag.UintVar(&cfg.AccrualWorkers, "workers", accrualWorkers, "Number of concurrent requests to Accrual")
	flag.UintVar(&cfg.AccrualInterval, "interval", accrualInterval, "Interval for generating requests to Accrual, sec")
	flag.UintVar(&cfg.AccrualShutdown, "accshutdown", accrualShutdown, "Accrual shutdown timeout, sec")

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
