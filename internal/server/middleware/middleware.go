package middleware

import (
	"net/http"
	"strings"

	"github.com/KryukovO/gophermart/internal/utils"
	"github.com/google/uuid"

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

func (mw *Manager) LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(echoCtx echo.Context) error {
		uuid := uuid.New()

		echoCtx.Set("uuid", uuid)

		mw.logger.Infof(
			"[%s] Request received with %s method: %s",
			uuid, echoCtx.Request().Method, echoCtx.Request().URL.Path,
		)

		return next(echoCtx)
	})
}

func (mw *Manager) GZipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(echoCtx echo.Context) error {
		uuid := echoCtx.Get("uuid")
		if uuid == nil {
			uuid = ""
		}

		contentEncoding := echoCtx.Request().Header.Get("Content-Encoding")
		sendsGZip := strings.Contains(contentEncoding, "gzip")

		if sendsGZip {
			reader, err := NewReader(echoCtx.Request().Body)
			if err != nil {
				mw.logger.Errorf("[%s] Something went wrong: %s", uuid, err)
				return echoCtx.NoContent(http.StatusInternalServerError)
			}

			defer reader.Close()

			echoCtx.Request().Body = reader
		}

		acceptEnc := echoCtx.Request().Header.Get("Accept-Encoding")
		supportGZip := strings.Contains(acceptEnc, "gzip")

		if supportGZip {
			writer := NewWriter(echoCtx.Response().Writer)
			echoCtx.Response().Writer = writer
		}

		return next(echoCtx)
	})
}

func (mw *Manager) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
