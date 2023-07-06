package gophermart

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	repoTimeout := time.Duration(cfg.RepositioryTimeout) * time.Second

	repoCtx, cancel := context.WithTimeout(context.Background(), repoTimeout)
	defer cancel()

	pg, err := postgres.NewPostgres(repoCtx, cfg.DSN, cfg.Migrations)
	if err != nil {
		return err
	}

	defer func() {
		pg.Close()

		logger.Info("Database connection closed")
	}()

	user := usecases.NewUserUseCase(pgrepo.NewUserRepo(pg), repoTimeout)
	order := usecases.NewOrderUseCase(pgrepo.NewOrderRepo(pg), repoTimeout)
	balance := usecases.NewBalanceUseCase(pgrepo.NewBalanceRepo(pg), repoTimeout)

	server, err := server.NewServer(
		cfg.Address, []byte(cfg.SecretKey),
		time.Duration(cfg.UserTokenTTL)*time.Minute,
		user, order, balance,
		logger,
	)
	if err != nil {
		return err
	}

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
		select {
		case <-groupCtx.Done():
			return nil
		case <-sigCtx.Done():
		}

		logger.Info("Stopping server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.ShutdownTimeout)*time.Second,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("Unable to gracefully stop the server: %s", err)
		} else {
			logger.Info("Server stopped gracefully")
		}

		return nil
	})

	return group.Wait()
}
