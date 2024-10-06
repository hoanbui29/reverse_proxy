package main

import (
	"hoanbui29/reverse_proxy/internal/logger"
	"net/http"
	"os"
)

type server struct {
	logger *logger.Logger
}

func main() {
	s := server{
		logger: logger.New(os.Stdout, logger.LevelDebug),
	}
	err := s.serve()
	if err != nil {
		s.logger.Fatal("server failed to start", map[string]string{"error": err.Error()})
	}
}

func (s server) serve() error {
	mux := routes()
	return http.ListenAndServe(":8080", mux)
}
