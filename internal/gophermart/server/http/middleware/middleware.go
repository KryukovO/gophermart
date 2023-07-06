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
	return echo.HandlerFunc(func(e echo.Context) error {
		uuid := uuid.New()

		e.Set("uuid", uuid)

		mw.logger.Infof(
			"[%s] Request received with %s method: %s",
			uuid, e.Request().Method, e.Request().URL.Path,
		)

		return next(e)
	})
}

func (mw *Manager) GZipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(e echo.Context) error {
		uuid := e.Get("uuid")
		if uuid == nil {
			uuid = ""
		}

		contentEncoding := e.Request().Header.Get("Content-Encoding")
		sendsGZip := strings.Contains(contentEncoding, "gzip")

		if sendsGZip {
			reader, err := NewReader(e.Request().Body)
			if err != nil {
				mw.logger.Errorf("[%s] Something went wrong: %s", uuid, err)

				return e.NoContent(http.StatusInternalServerError)
			}

			defer reader.Close()

			e.Request().Body = reader
		}

		acceptEnc := e.Request().Header.Get("Accept-Encoding")
		supportGZip := strings.Contains(acceptEnc, "gzip")

		if supportGZip {
			acceptTypes := [...]string{
				"application/json",
				"text/html",
			}
			writer := NewWriter(e.Response().Writer, acceptTypes[:])
			e.Response().Writer = writer
		}

		return next(e)
	})
}

func (mw *Manager) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(e echo.Context) error {
		tokenCookie, err := e.Cookie("token")
		if err != nil {
			return e.NoContent(http.StatusUnauthorized)
		}

		var userID int

		err = utils.ParseTokenString(&userID, tokenCookie.Value, mw.secret)
		if err != nil {
			return e.NoContent(http.StatusUnauthorized)
		}

		e.Set("userID", userID)

		return next(e)
	})
}
