package main

import (
	"flag"

	"github.com/KryukovO/loyalty/internal/server"
	"github.com/KryukovO/loyalty/internal/server/config"

	log "github.com/sirupsen/logrus"
)

const (
	address        = "127.0.0.1:8080"
	dsn            = ""
	accrualAddress = ""
)

func main() {
	cfg := new(config.Config)

	flag.StringVar(&cfg.Address, "a", address, "Server address")
	flag.StringVar(&cfg.DSN, "d", dsn, "Data source name")
	flag.StringVar(&cfg.AccrualAddress, "r", accrualAddress, "Accrual system address")
	flag.Parse()

	logger := log.New()
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05 Z07:00",
	})

	srv := server.NewServer(cfg, logger)
	if err := srv.Run(); err != nil {
		logger.Fatalf("Server error: %s. Exit(1)", err)
	}
}
