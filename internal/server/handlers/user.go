package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/KryukovO/gophermart/internal/entities"
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/KryukovO/gophermart/internal/utils"
	"github.com/labstack/echo"

	log "github.com/sirupsen/logrus"
)

type UserController struct {
	user          usecases.User
	secret        []byte
	tokenLifetime time.Duration
	logger        *log.Logger
}

func NewUserController(
	user usecases.User,
	secret []byte, tokenLifetime time.Duration,
	logger *log.Logger,
) (*UserController, error) {
	if user == nil {
		return nil, ErrUseCaseIsNil
	}

	controllerLogger := log.StandardLogger()
	if logger != nil {
		controllerLogger = logger
	}

	return &UserController{
		user:          user,
		secret:        secret,
		tokenLifetime: tokenLifetime,
		logger:        controllerLogger,
	}, nil
}

func (c *UserController) MapHandlers(router *echo.Router) error {
	if router == nil {
		return ErrRouterIsNil
	}

	router.Add(http.MethodPost, "/api/user/register", c.registerHandler)
	router.Add(http.MethodPost, "/api/user/login", c.loginHandler)

	return nil
}

func (c *UserController) registerHandler(e echo.Context) error {
	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		c.logger.Errorf("Something went wrong: %s", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	var user entities.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	err = c.user.Register(e.Request().Context(), &user, c.secret)
	if err != nil {
		if errors.Is(err, entities.ErrUserAlreadyExists) {
			return e.NoContent(http.StatusConflict)
		}

		c.logger.Errorf("Something went wrong: %s", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	return e.NoContent(http.StatusOK)
}

func (c *UserController) loginHandler(e echo.Context) error {
	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		c.logger.Errorf("Something went wrong: %s", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	var user entities.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		return e.NoContent(http.StatusBadRequest)
	}

	err = c.user.Login(e.Request().Context(), &user, c.secret)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidLoginPassword) {
			return e.NoContent(http.StatusUnauthorized)
		}

		c.logger.Errorf("Something went wrong: %s", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	tokenString, err := utils.BuildJSWTString(c.secret, c.tokenLifetime, user.ID)
	if err != nil {
		c.logger.Errorf("Something went wrong: %s", err)
		return e.NoContent(http.StatusInternalServerError)
	}

	e.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
	})

	return e.NoContent(http.StatusOK)
}
