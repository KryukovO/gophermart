package server

import (
	"context"
	"errors"

	"github.com/KryukovO/gophermart/internal/server/handlers"
	"github.com/KryukovO/gophermart/internal/usecases"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

var ErrUseCaseIsNil = errors.New("usecase is nil")

type Server struct {
	address    string
	httpServer *echo.Echo
	logger     *log.Logger
}

func NewServer(
	address string,
	user usecases.User, order usecases.Order, balance usecases.Balance,
	logger *log.Logger,
) (*Server, error) {
	if user == nil {
		return nil, ErrUseCaseIsNil
	}

	serverLogger := log.StandardLogger()
	if logger != nil {
		serverLogger = logger
	}

	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.HidePort = true

	err := handlers.SetHandlers(httpServer.Router(), user, order, balance, logger)
	if err != nil {
		return nil, err
	}

	return &Server{
		address:    address,
		httpServer: httpServer,
		logger:     serverLogger,
	}, nil
}

func (s *Server) Run() error {
	return s.httpServer.Start(s.address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
