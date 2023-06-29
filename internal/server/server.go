package server

import (
	"github.com/KryukovO/loyalty/internal/server/config"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	cfg    *config.Config
	logger *log.Logger
}

func NewServer(cfg *config.Config, logger *log.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	s.logger.Infof("Run server at %s", s.cfg.Address)

	return nil
}
