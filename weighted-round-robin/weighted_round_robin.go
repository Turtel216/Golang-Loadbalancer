package weighted_round_robin

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

type Server interface {
	Address() string
	Serve(rw http.ResponseWriter, req *http.Request)
}

type simpleServer struct {
	addr   string
	weight int
	proxy  *httputil.ReverseProxy
}

// Creates a new instance of the simpleServer struct
func NewSimpleServer(addr string, weight int) simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	return simpleServer{
		addr:   addr,
		weight: weight,
		proxy:  httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Loadbalancer struct {
	Port          string
	mu            sync.Mutex
	servers       []simpleServer
	current       int
	total_weights int
}

// Creates and returns a new loadbalancer instance
func NewLoadBalancer(port string, servers []simpleServer, total_weights int) *Loadbalancer {
	loadbalancer := &Loadbalancer{
		Port:          port,
		servers:       servers,
		total_weights: total_weights,
	}

	for _, server := range servers {
		loadbalancer.total_weights += server.weight
	}

	return loadbalancer
}

// Returns the adress of the simple server instance
func (s *simpleServer) Address() string { return s.addr }

// Serves the through the reverse proxy
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// returns the server selected by the weighted round-robin scheduler
func (loadbalancer *Loadbalancer) getNextAvailableServer() (simpleServer, error) {
	loadbalancer.mu.Lock()
	defer loadbalancer.mu.Unlock()

	if loadbalancer.current == -1 {
		loadbalancer.current = rand.Intn(loadbalancer.total_weights)
	}

	for i := 0; i < loadbalancer.total_weights; i++ {
		loadbalancer.current = (loadbalancer.current + 1) + loadbalancer.total_weights
		loadbalancer.servers[loadbalancer.current].weight += loadbalancer.servers[loadbalancer.current].weight

		if loadbalancer.servers[loadbalancer.current].weight >= loadbalancer.total_weights {
			loadbalancer.servers[loadbalancer.current].weight -= loadbalancer.total_weights
			return loadbalancer.servers[loadbalancer.current], nil
		}
	}

	return simpleServer{}, errors.New("Couldn't find next avaible server")
}

// Forwards the request to the server returned by the getNextAvailableServer method
func (loadbalancer *Loadbalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	target, err := loadbalancer.getNextAvailableServer()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Forwarding request to adress %q\n", target.Address())

	target.Serve(rw, req)
}
