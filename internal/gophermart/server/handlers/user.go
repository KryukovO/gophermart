package handlers

import (
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type UserController struct {
	user   usecases.User
	logger *log.Logger
}

func NewUserController(user usecases.User, logger *log.Logger) (*UserController, error) {
	if user == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &UserController{
		user:   user,
		logger: controllerLogger,
	}, nil
}

func (controller *UserController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	return nil
}
