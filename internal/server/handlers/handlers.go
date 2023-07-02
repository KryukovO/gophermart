package handlers

import (
	"errors"

	"github.com/KryukovO/gophermart/internal/server/middleware"
	"github.com/KryukovO/gophermart/internal/usecases"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var (
	ErrUseCaseIsNil = errors.New("usecase is nil")
	ErrRouterIsNil  = errors.New("router is nil")
)

func SetHandlers(
	router *echo.Router, secret []byte,
	user usecases.User, order usecases.Order, balance usecases.Balance,
	logger *log.Logger,
) error {
	if router == nil {
		return ErrRouterIsNil
	}

	mwManager := middleware.NewManager(secret, logger)

	userController, err := NewUserController(user, logger)
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

	err = userController.MapHandlers(router)
	if err != nil {
		return err
	}

	err = orderController.MapHandlers(router)
	if err != nil {
		return err
	}

	err = balanceController.MapHandlers(router)
	if err != nil {
		return err
	}

	return nil
}
