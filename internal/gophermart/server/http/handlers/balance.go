package handlers

import (
	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type BalanceController struct {
	balance usecases.Balance
	mw      *middleware.Manager
	logger  *log.Logger
}

func NewBalanceController(
	balance usecases.Balance, mwManager *middleware.Manager, logger *log.Logger,
) (*BalanceController, error) {
	if balance == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &BalanceController{
		balance: balance,
		mw:      mwManager,
		logger:  controllerLogger,
	}, nil
}

func (controller *BalanceController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	return nil
}
