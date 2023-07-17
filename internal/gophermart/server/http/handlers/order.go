package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"

	"github.com/labstack/echo/v4"
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

func (c *OrderController) MapHandlers(group *echo.Group) error {
	if group == nil {
		return ErrGroupIsNil
	}

	group.Add(http.MethodPost, "/user/orders", c.mw.AuthenticationMiddleware(c.addOrderHandler))
	group.Add(http.MethodGet, "/user/orders", c.mw.AuthenticationMiddleware(c.ordersHandler))

	return nil
}

// @Summary Add new order
// @ID add-order
// @Accept json
// @Success 200
// @Success 202
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/user/orders [post]
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

	c.logger.Debugf("[%s] Request body: %s", uuid, string(body))

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

// @Summary Get uploaded orders
// @ID orders
// @Accept plain
// @Produce json
// @Success 200 {object} entities.Order
// @Success 204
// @Failure 401 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/user/orders [get]
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

	return e.JSON(http.StatusOK, orders)
}
