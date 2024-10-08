package main

import (
	"hoanbui29/reverse_proxy/internal/config"
	"hoanbui29/reverse_proxy/internal/logger"
	"sync"
)

type serverPool interface {
	AddServer(*server)
	Next() *server
}

type server struct {
	address         string
	liveConnections int
	mu              *sync.Mutex
}

type loadBalancer struct {
	servers []*server             //Shared servers between resources
	poolMap map[string]serverPool //Concrete servers by resource
}

type serverGateway struct {
	logger *logger.Logger
	config config.Config
	lb     *loadBalancer
}

func (s server) Add(amount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.liveConnections += amount
}

func (s *serverGateway) loadConfig() {
	cfg, err := config.Load()
	s.handleFatal(err)
	s.config = cfg
}

func (s *serverGateway) loadLBConfig() {
	lb := loadBalancer{
		servers: make([]*server, 0),
		poolMap: make(map[string]serverPool),
	}
	checkMap := make(map[string]*server)
	for _, resource := range s.config.Resources {
		pool, err := getServerPool(&resource)
		if err != nil {
			panic(err)
		}
		lb.poolMap[resource.Prefix] = pool
		for _, dest := range resource.Destinations {
			var srv *server
			if s, ok := checkMap[dest]; ok {
				srv = s
			} else {
				srv = &server{
					address: dest,
					mu:      &sync.Mutex{},
				}
				checkMap[dest] = srv
				lb.servers = append(lb.servers, srv)
			}
			lb.poolMap[resource.Prefix].AddServer(srv)
		}
	}
	s.lb = &lb
}
