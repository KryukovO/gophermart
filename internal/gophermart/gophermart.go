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
	"github.com/KryukovO/gophermart/internal/gophermart/server"
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/KryukovO/gophermart/internal/usecases/repository/pgrepo"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

func Run(cfg *config.Config, logger *log.Logger) error {
	logger.Info("Connecting to the repository...")

	repo, err := pgrepo.NewPgRepo(cfg.DSN, cfg.Migrations)
	if err != nil {
		return err
	}

	defer func() {
		repo.Close()

		logger.Info("Repository closed")
	}()

	user := usecases.NewUserUseCase(repo)

	server, err := server.NewServer(cfg.Address, user, logger)
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
