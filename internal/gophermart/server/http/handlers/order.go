package handlers

import (
	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type OrderController struct {
	order  usecases.Order
	mw     *middleware.Manager
	logger *log.Logger
}

func NewOrderController(
	order usecases.Order, mwManager *middleware.Manager, logger *log.Logger,
) (*OrderController, error) {
	if order == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &OrderController{
		order:  order,
		mw:     mwManager,
		logger: controllerLogger,
	}, nil
}

func (controller *OrderController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	return nil
}
