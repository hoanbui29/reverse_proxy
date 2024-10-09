package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"hoanbui29/reverse_proxy/internal/config"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

type rateClient struct {
	*rate.Limiter
	lastConnected time.Time
}

type rateClientMap struct {
	mu      *sync.Mutex
	clients map[string]*rateClient
}

func newRateClientMap() *rateClientMap {
	return &rateClientMap{
		mu:      &sync.Mutex{},
		clients: make(map[string]*rateClient),
	}
}

func newRateClient(cfg config.Config) *rateClient {
	return &rateClient{
		Limiter:       rate.NewLimiter(rate.Limit(cfg.Server.RateLimit.RatePerSecond), cfg.Server.RateLimit.Burst),
		lastConnected: time.Now(),
	}
}

func (rcm *rateClientMap) getRateLimiter(cfg config.Config, clientIP string) *rateClient {
	rcm.mu.Lock()
	defer rcm.mu.Unlock()
	seenTime := time.Now()
	if v, ok := rcm.clients[clientIP]; ok {
		v.lastConnected = seenTime
		return v
	}

	rcm.clients[clientIP] = newRateClient(cfg)
	return rcm.clients[clientIP]
}

func (s serverGateway) getResourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range s.config.Resources {
			if strings.HasPrefix(r.URL.Path, v.Prefix) {
				r = setResource(r, &v)
				next.ServeHTTP(w, r)
				return
			}
		}
		s.notFound(w, r)
		return
	})
}

func (s serverGateway) checkAllowedMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := getResource(r)
		if !slices.Contains(resource.Methods, r.Method) {
			s.methodNotAllowed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s serverGateway) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error(err.(error), nil)
				s.internalServerError(w, r, err.(error))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *serverGateway) rateLimiter(next http.Handler) http.Handler {
	rcm := newRateClientMap()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error(fmt.Errorf("panic: %v", err), nil)
			}
		}()

		for {
			rcm.mu.Lock()
			for k, v := range rcm.clients {
				if time.Since(v.lastConnected) > 5*time.Minute {
					delete(rcm.clients, k)
				}
			}

			rcm.mu.Unlock()
			time.Sleep(1 * time.Minute)
		}

	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, err := getClientIP(r)

		if err != nil {
			s.internalServerError(w, r, err)
			return
		}
		lm := rcm.getRateLimiter(s.config, clientIP)
		if !lm.Allow() {
			s.rateLimitExceeded(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
