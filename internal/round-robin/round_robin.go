package round_robin

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	lb "github.com/Turtel216/Golang-Loadbalancer/internal"
)

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

// Creates and returns a new instance of the simpleServer struct
func NewSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Loadbalancer struct {
	Port            string
	roundRobinCount int

	mu      sync.Mutex
	servers []lb.Server //Interface
}

// Creates a new loadbalancer instance
func NewLoadBalancer(port string, servers []lb.Server) *Loadbalancer {
	return &Loadbalancer{
		Port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

// Returns the server selected by the round-robin scheduler
func (loadbalancer *Loadbalancer) getNextAvailableServer() (server lb.Server) {
	loadbalancer.mu.Lock()
	defer loadbalancer.mu.Unlock()

	server = loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]

	for !server.IsAlive() {
		loadbalancer.roundRobinCount++
		server = loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]
	}

	loadbalancer.roundRobinCount++
	return
}

// Returns the adress of the simple server instance
func (s *simpleServer) Address() string { return s.addr }

// Ensures that the simpleServer instsance is running
func (s *simpleServer) IsAlive() bool { return true }

// Serves the through the reverse proxy
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// Forwards the request to the server returned by the getNextAvailableServer method
func (loadbalancer *Loadbalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	target := loadbalancer.getNextAvailableServer()
	fmt.Printf("Forwarding request to adress %q\n", target.Address())

	target.Serve(rw, req)
}
