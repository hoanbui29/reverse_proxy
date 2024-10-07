package main

import (
	"fmt"
	"hoanbui29/reverse_proxy/internal/config"
	"hoanbui29/reverse_proxy/internal/logger"
	"net/http"
	"os"
)

type serverGateway struct {
	logger *logger.Logger
	config config.Config
}

func main() {
	s := serverGateway{
		logger: logger.New(os.Stdout, logger.LevelDebug),
	}
	cfg, err := config.Load()
	s.handleFatal(err, "failed to load config")
	s.config = cfg
	s.handleFatal(s.serve(), "failed to start server")
}

func (s serverGateway) serve() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port), s.proxy())
}
