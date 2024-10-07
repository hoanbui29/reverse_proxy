package main

import (
	"errors"
	"hoanbui29/reverse_proxy/internal/config"
	"sync"
)

var ErrStrategyNotSupported = errors.New("strategy not supported")

type IResourceFactory interface {
	Next() string
}

func GetResourceFactory(r *config.Resource) (IResourceFactory, error) {
	switch r.Strategy {
	case config.StrategyRoundRobin:
		return &RoundRobinResource{
			hosts: r.Destinations,
			index: 0,
			mu:    &sync.Mutex{},
		}, nil
	default:
		return nil, ErrStrategyNotSupported
	}
}

type RoundRobinResource struct {
	hosts []string
	index int
	mu    *sync.Mutex
}

func (r *RoundRobinResource) Next() string {
	// Implement the logic to return the next host in the list
	r.mu.Lock()
	defer r.mu.Unlock()

	result := r.hosts[r.index]

	r.index = (r.index + 1) % len(r.hosts)

	return result
}
