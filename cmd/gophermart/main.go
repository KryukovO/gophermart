package main

import (
	"github.com/KryukovO/gophermart/internal/gophermart"
	"github.com/KryukovO/gophermart/internal/gophermart/config"
	"github.com/spf13/pflag"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfig()
	helpFlag := false

	pflag.BoolVarP(&helpFlag, "help", "h", false, "Shows gophermart usage")

	pflag.StringVarP(&cfg.Address, "address", "a", cfg.Address, "Address to run HTTP server")
	pflag.StringVarP(&cfg.DSN, "dsn", "d", cfg.DSN, "URI to database")
	pflag.StringVarP(&cfg.AccrualAddress, "accrual", "r", cfg.AccrualAddress, "Accrual system address")

	pflag.StringVar(&cfg.SecretKey, "secret", cfg.SecretKey, "Authorization token encryption key")
	pflag.DurationVar(&cfg.UserTokenTTL, "userttl", cfg.UserTokenTTL, "User token lifetime")
	pflag.DurationVar(&cfg.ShutdownTimeout, "shutdown", cfg.ShutdownTimeout, "Server shutdown timeout")
	pflag.DurationVar(&cfg.RepositioryTimeout, "timeout", cfg.RepositioryTimeout, "Repository connection timeout")
	pflag.StringVar(&cfg.Migrations, "migrations", cfg.Migrations, "Directory of database migration files")
	pflag.UintVar(&cfg.AccrualWorkers, "workers", cfg.AccrualWorkers, "Number of concurrent requests to Accrual")
	pflag.DurationVar(&cfg.AccrualInterval, "interval", cfg.AccrualInterval, "Interval for generating requests to Accrual")
	pflag.DurationVar(&cfg.AccrualShutdown, "accshutdown", cfg.AccrualShutdown, "Accrual connector shutdown timeout")

	pflag.Parse()

	if helpFlag {
		pflag.Usage()

		return
	}

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05 Z07:00",
	})

	if err := gophermart.Run(cfg, logger); err != nil {
		logger.Fatalf("Gophermart service error: %s. Exit(1)", err)
	}
}
