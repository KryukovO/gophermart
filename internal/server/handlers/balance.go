package handlers

import (
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type BalanceController struct {
	balance usecases.Balance
	logger  *log.Logger
}

func NewBalanceController(balance usecases.Balance, logger *log.Logger) (*BalanceController, error) {
	if balance == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &BalanceController{
		balance: balance,
		logger:  controllerLogger,
	}, nil
}

func (controller *BalanceController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	return nil
}
