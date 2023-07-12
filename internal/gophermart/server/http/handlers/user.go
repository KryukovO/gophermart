package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
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

	router.Add(http.MethodPost, "/api/user/register", c.userRequestHandler(c.user.Register))
	router.Add(http.MethodPost, "/api/user/login", c.userRequestHandler(c.user.Login))

	return nil
}

func (c *UserController) userRequestHandler(
	userFunc func(ctx context.Context, user *entities.User, secret []byte) error,
) func(echo.Context) error {
	return func(e echo.Context) error {
		uuid := e.Get("uuid")
		if uuid == nil {
			uuid = ""
		}

		body, err := io.ReadAll(e.Request().Body)
		if err != nil {
			c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

			return e.NoContent(http.StatusInternalServerError)
		}

		var user entities.User

		err = json.Unmarshal(body, &user)
		if err != nil {
			return e.NoContent(http.StatusBadRequest)
		}

		c.logger.Debugf("[%s] Request body: %+v", uuid, user)

		err = userFunc(e.Request().Context(), &user, c.secret)
		if err != nil {
			if errors.Is(err, entities.ErrUserAlreadyExists) {
				return e.NoContent(http.StatusConflict)
			}

			if errors.Is(err, entities.ErrInvalidLoginPassword) {
				return e.NoContent(http.StatusUnauthorized)
			}

			c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

			return e.NoContent(http.StatusInternalServerError)
		}

		tokenString, err := utils.BuildJSWTString(c.secret, c.tokenLifetime, user.ID)
		if err != nil {
			c.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

			return e.NoContent(http.StatusInternalServerError)
		}

		e.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		})

		return e.NoContent(http.StatusOK)
	}
}
