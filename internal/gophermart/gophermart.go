package gophermart

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KryukovO/gophermart/internal/gophermart/accrualconnector"
	"github.com/KryukovO/gophermart/internal/gophermart/config"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/pgrepo"
	server "github.com/KryukovO/gophermart/internal/gophermart/server/http"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/KryukovO/gophermart/internal/postgres"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *config.Config, logger *log.Logger) error {
	logger.Infof("Connect to the database: %s", cfg.DSN)

	repoCtx, cancel := context.WithTimeout(context.Background(), cfg.RepositioryTimeout)
	defer cancel()

	pg, err := postgres.NewPostgres(repoCtx, cfg.DSN, cfg.Migrations)
	if err != nil {
		return err
	}

	logger.Info("Database connection established")

	defer func() {
		pg.Close()

		logger.Info("Database connection closed")
	}()

	user := usecases.NewUserUseCase(pgrepo.NewUserRepo(pg), cfg.RepositioryTimeout)
	order := usecases.NewOrderUseCase(pgrepo.NewOrderRepo(pg), cfg.RepositioryTimeout)
	balance := usecases.NewBalanceUseCase(pgrepo.NewBalanceRepo(pg), cfg.RepositioryTimeout)

	server, err := server.NewServer(
		cfg.Address, []byte(cfg.SecretKey),
		cfg.UserTokenTTL,
		user, order, balance,
		logger,
	)
	if err != nil {
		return err
	}

	accrualConnector := accrualconnector.NewAccrualConnector(
		cfg.AccrualAddress, cfg.AccrualWorkers, cfg.AccrualInterval,
		order, balance,
		logger,
	)

	sigCtx, sigCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer sigCancel()

	group, groupCtx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		logger.Infof("Run server at %s", cfg.Address)

		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	group.Go(func() error {
		logger.Infof("Run accrual connector: workers: %d, interval: %s", cfg.AccrualWorkers, cfg.AccrualInterval)

		accrualConnector.Run(groupCtx)

		logger.Info("Accrual connector stopped")

		return nil
	})

	group.Go(func() error {
		select {
		case <-groupCtx.Done():
			return nil
		case <-sigCtx.Done():
			logger.Info("Shutdown signal received")
		}

		logger.Info("Stopping server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			cfg.ShutdownTimeout,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("Unable to gracefully stop the server: %s", err)
		} else {
			logger.Info("Server stopped gracefully")
		}

		accrualCtx, accrualCancel := context.WithTimeout(
			context.Background(),
			cfg.AccrualShutdown,
		)
		defer accrualCancel()

		accrualConnector.Shutdown(accrualCtx)

		return nil
	})

	return group.Wait()
}
