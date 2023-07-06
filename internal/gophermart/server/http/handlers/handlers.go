package handlers

import (
	"errors"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUseCaseIsNil = errors.New("usecase is nil")
	ErrServerIsNil  = errors.New("server instance is nil")
	ErrRouterIsNil  = errors.New("router is nil")
)

func SetHandlers(
	server *echo.Echo,
	secret []byte, tokenLifetime time.Duration,
	user usecases.User, order usecases.Order, balance usecases.Balance,
	logger *log.Logger,
) error {
	if server == nil {
		return ErrRouterIsNil
	}

	mwManager := middleware.NewManager(secret, logger)

	userController, err := NewUserController(user, secret, tokenLifetime, logger)
	if err != nil {
		return err
	}

	orderController, err := NewOrderController(order, mwManager, logger)
	if err != nil {
		return err
	}

	balanceController, err := NewBalanceController(balance, mwManager, logger)
	if err != nil {
		return err
	}

	err = userController.MapHandlers(server.Router())
	if err != nil {
		return err
	}

	err = orderController.MapHandlers(server.Router())
	if err != nil {
		return err
	}

	err = balanceController.MapHandlers(server.Router())
	if err != nil {
		return err
	}

	server.Use(
		mwManager.LoggingMiddleware,
		mwManager.GZipMiddleware,
	)

	return nil
}
