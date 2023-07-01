package handlers

import (
	"errors"

	"github.com/KryukovO/gophermart/internal/usecases"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var (
	ErrServerIsNil  = errors.New("server instance is nil")
	ErrUseCaseIsNil = errors.New("usecase is nil")
	ErrRouterIsNil  = errors.New("router is nil")
)

func SetHandlers(server *echo.Echo, user usecases.User, logger *log.Logger) error {
	if server == nil {
		return ErrServerIsNil
	}

	userController, err := NewUserController(user, logger)
	if err != nil {
		return err
	}

	return userController.MapHandlers(server.Router())
}
