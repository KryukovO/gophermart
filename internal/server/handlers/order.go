package handlers

import (
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type OrderController struct {
	order  usecases.Order
	logger *log.Logger
}

func NewOrderController(order usecases.Order, logger *log.Logger) (*OrderController, error) {
	if order == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &OrderController{
		order:  order,
		logger: controllerLogger,
	}, nil
}

func (controller *OrderController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	return nil
}
