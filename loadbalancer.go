package main

import (
	"fmt"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	servers []string
	index   int
	mu      sync.Mutex
	count   int
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		index:   0,
		count:   0,
	}
}

func (lb *LoadBalancer) NextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return ""
	}

	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.count++
	if mod := lb.count % 1000; mod == 0 {
		fmt.Printf("LB req count: %d\n", lb.count)
	}

	server := lb.NextServer()

	if server == "" {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	proxy := NewProxy(server)
	proxy.ServeHTTP(w, r)
}
