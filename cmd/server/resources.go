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

type powerOfTwoChoicesPool struct {
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
	case config.StrategyLeastConnections:
		return &leastConnectionsPool{
			servers: []*server{},
			mu:      &sync.Mutex{},
		}, nil
	default:
		return nil, ErrStrategyNotSupported
	}
}

func (p *roundRobinPool) Next() *server {
	// Implement the logic to return the next host in the list
	p.mu.Lock()
	defer p.mu.Unlock()

	srv := p.servers[p.index]

	p.index = (p.index + 1) % len(p.servers)

	return srv
}

func (p *roundRobinPool) AddServer(s *server) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, s)
}

func (p *leastConnectionsPool) Next() *server {
	if len(p.servers) == 0 {
		panic("no servers available")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	min := p.servers[0].liveConnections
	minIndex := 0

	for i := 1; i < len(p.servers); i++ {
		if p.servers[i].liveConnections < min {
			min = p.servers[i].liveConnections
			minIndex = i
		}
	}

	return p.servers[minIndex]
}

func (p *leastConnectionsPool) AddServer(s *server) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, s)
}

func (p *powerOfTwoChoicesPool) Next() *server {
	if len(p.servers) == 0 {
		panic("no servers available")
	}

	if len(p.servers) == 1 {
		return p.servers[0]
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	srvs := randomPick(p.servers, 2)

	if srvs[0].CountLiveCons() < srvs[1].CountLiveCons() {
		return srvs[0]
	}

	return srvs[1]
}

func (p *powerOfTwoChoicesPool) AddServer(s *server) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, s)
}
