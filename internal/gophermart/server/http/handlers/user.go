package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/KryukovO/gophermart/internal/utils"
	"github.com/labstack/echo/v4"

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

func (c *UserController) MapHandlers(group *echo.Group) error {
	if group == nil {
		return ErrGroupIsNil
	}

	group.Add(http.MethodPost, "/user/register", c.registerHandler)
	group.Add(http.MethodPost, "/user/login", c.loginHandler)

	return nil
}

// @Summary       User registration
// @Description   User registration by login and password.
// @Tags          Gophermart HTTP API
// @Accept        json
// @Param         user   body       entities.User   true   "User login and password."
// @Success       200
// @Failure       400    {object}   echo.HTTPError
// @Failure       409    {object}   echo.HTTPError
// @Failure       500    {object}   echo.HTTPError
// @Router        /api/user/register [post]
func (c *UserController) registerHandler(e echo.Context) error {
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

	if user.Login == "" || user.Password == "" {
		return e.NoContent(http.StatusBadRequest)
	}

	c.logger.Debugf("[%s] Request body: %+v", uuid, user)

	err = c.user.Register(e.Request().Context(), &user, c.secret)
	if err != nil {
		if errors.Is(err, entities.ErrUserAlreadyExists) {
			return e.NoContent(http.StatusConflict)
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

// @Summary       User authorization
// @Description   User authorization by login and password.
// @Tags          Gophermart HTTP API
// @Accept        json
// @Param         user   body       entities.User   true   "User login and password."
// @Success       200
// @Failure       400    {object}   echo.HTTPError
// @Failure       401    {object}   echo.HTTPError
// @Failure       500    {object}   echo.HTTPError
// @Router        /api/user/login [post]
func (c *UserController) loginHandler(e echo.Context) error {
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

	if user.Login == "" || user.Password == "" {
		return e.NoContent(http.StatusBadRequest)
	}

	c.logger.Debugf("[%s] Request body: %+v", uuid, user)

	err = c.user.Login(e.Request().Context(), &user, c.secret)
	if err != nil {
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
