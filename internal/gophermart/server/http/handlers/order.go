package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
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

func (c *OrderController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	router.Add(http.MethodPost, "/api/user/orders", c.mw.AuthenticationMiddleware(c.addOrderHandler))
	router.Add(http.MethodGet, "/api/user/orders", c.mw.AuthenticationMiddleware(c.ordersHandler))

	return nil
}

func (c *OrderController) addOrderHandler(e echo.Context) error {
	uuid := e.Get("uuid")
	if uuid == nil {
		uuid = ""
	}

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	userID := e.Get("userID")

	user, ok := userID.(int64)
	if !ok {
		return e.NoContent(http.StatusUnauthorized)
	}

	order := entities.NewOrder(string(body), user)

	err = c.order.AddOrder(e.Request().Context(), order)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidOrderNumber) {
			return e.NoContent(http.StatusUnprocessableEntity)
		}

		if errors.Is(err, entities.ErrOrderAlreadyAdded) {
			return e.NoContent(http.StatusOK)
		}

		if errors.Is(err, entities.ErrOrderAddedByOther) {
			return e.NoContent(http.StatusConflict)
		}

		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	return e.NoContent(http.StatusAccepted)
}

func (c *OrderController) ordersHandler(e echo.Context) error {
	uuid := e.Get("uuid")
	if uuid == nil {
		uuid = ""
	}

	userID := e.Get("userID")

	user, ok := userID.(int64)
	if !ok {
		return e.NoContent(http.StatusUnauthorized)
	}

	orders, err := c.order.Orders(e.Request().Context(), user)
	if err != nil {
		c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	if len(orders) == 0 {
		return e.NoContent(http.StatusNoContent)
	}

	return e.NoContent(http.StatusOK)
}