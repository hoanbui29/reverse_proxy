package main

import (
	"errors"
	"hoanbui29/reverse_proxy/internal/config"
	"sync"
)

var ErrStrategyNotSupported = errors.New("strategy not supported")

type roundRobinPool struct {
	servers []*server
	index   int
	mu      *sync.Mutex
}

type leastConnectionsPool struct {
	servers []*server
	mu      *sync.Mutex
}

func getServerPool(r *config.Resource) (serverPool, error) {
	switch r.Strategy {
	case config.StrategyRoundRobin:
		return &roundRobinPool{
			servers: []*server{},
			index:   0,
			mu:      &sync.Mutex{},
		}, nil
	default:
		return nil, ErrStrategyNotSupported
	}
}

func (r *roundRobinPool) Next() *server {
	// Implement the logic to return the next host in the list
	r.mu.Lock()
	defer r.mu.Unlock()

	srv := r.servers[r.index]

	r.index = (r.index + 1) % len(r.servers)

	return srv
}

func (r *roundRobinPool) AddServer(s *server) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.servers = append(r.servers, s)
}
