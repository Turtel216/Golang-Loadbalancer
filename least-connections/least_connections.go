package least_connections

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

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
	Port        string
	Connections map[Server]int

	mu      sync.Mutex
	servers []Server //Interface
}

// Creates a new loadbalancer instance
func NewLoadBalancer(port string, servers []Server) *Loadbalancer {
	return &Loadbalancer{
		Port:        port,
		Connections: make(map[Server]int),
		servers:     servers,
	}
}

// Returns the server with the least connections
func (loadbalancer *Loadbalancer) getNextAvailableServer() (server Server) {
	loadbalancer.mu.Lock()
	defer loadbalancer.mu.Unlock()

	// Find server with least connections
	min_conn := 0

	for serv, num_conn := range loadbalancer.Connections {
		if num_conn < min_conn {
			min_conn = num_conn
			server = serv
		}
	}

	if server != nil {
		loadbalancer.Connections[server] = min_conn + 1
	}

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
