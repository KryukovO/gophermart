package handlers

import (
	"errors"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	ErrUseCaseIsNil = errors.New("usecase is nil")
	ErrServerIsNil  = errors.New("server instance is nil")
	ErrGroupIsNil   = errors.New("rout group is nil")
)

func SetHandlers(
	server *echo.Echo,
	secret []byte, tokenLifetime time.Duration,
	user usecases.User, order usecases.Order, balance usecases.Balance,
	logger *log.Logger,
) error {
	if server == nil {
		return ErrServerIsNil
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

	group := server.Group("/api")
	group.Use(
		mwManager.LoggingMiddleware,
		mwManager.GZipMiddleware,
	)

	err = userController.MapHandlers(group)
	if err != nil {
		return err
	}

	err = orderController.MapHandlers(group)
	if err != nil {
		return err
	}

	err = balanceController.MapHandlers(group)
	if err != nil {
		return err
	}

	server.GET("/swagger/*", echoSwagger.WrapHandler)

	return nil
}
