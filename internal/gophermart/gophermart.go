package gophermart

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KryukovO/gophermart/internal/config"
	"github.com/KryukovO/gophermart/internal/server"
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/KryukovO/gophermart/internal/usecases/repository/pgrepo"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *config.Config, logger *log.Logger) error {
	logger.Info("Connecting to the repository...")

	repoTimeout := time.Duration(cfg.RepositioryTimeout) * time.Second

	repoCtx, cancel := context.WithTimeout(context.Background(), repoTimeout)
	defer cancel()

	repo, err := pgrepo.NewPgRepo(repoCtx, cfg.DSN, cfg.Migrations)
	if err != nil {
		return err
	}

	defer func() {
		repo.Close()

		logger.Info("Repository closed")
	}()

	user := usecases.NewUserUseCase(repo, repoTimeout)
	order := usecases.NewOrderUseCase(repo, repoTimeout)
	balance := usecases.NewBalanceUseCase(repo, repoTimeout)

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
