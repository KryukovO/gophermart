package middleware

import (
	"net/http"

	"github.com/KryukovO/gophermart/internal/utils"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	secret []byte
	logger *log.Logger
}

func NewManager(secret []byte, logger *log.Logger) *Manager {
	middlewareLogger := log.StandardLogger()
	if logger != nil {
		middlewareLogger = logger
	}

	return &Manager{
		secret: secret,
		logger: middlewareLogger,
	}
}

func (mw *Manager) Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(echoCtx echo.Context) error {
		tokenCookie, err := echoCtx.Cookie("token")
		if err != nil {
			return echoCtx.NoContent(http.StatusUnauthorized)
		}

		var userID int

		err = utils.ParseTokenString(&userID, tokenCookie.Value, mw.secret)
		if err != nil {
			return echoCtx.NoContent(http.StatusUnauthorized)
		}

		echoCtx.Set("userID", userID)

		return next(echoCtx)
	})
}
