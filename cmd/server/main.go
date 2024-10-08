package main

import (
	"fmt"
	"hoanbui29/reverse_proxy/internal/logger"
	"net/http"
	"os"
)

func main() {
	s := serverGateway{
		logger: logger.New(os.Stdout, logger.LevelDebug),
	}
	s.loadConfig()
	s.loadLBConfig()
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port), s.getResourceMiddleware(s.checkAllowedMethods(http.HandlerFunc(s.loadBalancerHandler))))
}
