package source_ip_hash

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

// Creates a new instance of the simpleServer struct
func NewSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Loadbalancer struct {
	Port    string
	mu      sync.Mutex
	servers []SimpleServer
}

// Creates and returns a new loadbalancer instance
func NewLoadBalancer(port string, servers []SimpleServer) *Loadbalancer {
	loadbalancer := &Loadbalancer{
		Port:    port,
		servers: servers,
	}

	return loadbalancer
}

// Returns the server selected by the source ip hash algorithm
func (loadbalancer *Loadbalancer) getNextAvailableServer(req *http.Request) (SimpleServer, error) {
	loadbalancer.mu.Lock()
	defer loadbalancer.mu.Unlock()
	// get the source ip
	request_ip := req.Header.Get("X-Forwarded-For")

	// calculate the hash of the source ip
	ip_hash := hash(request_ip)

	// Map the hash value to the server index
	server_index := int(ip_hash) % len(loadbalancer.servers)

	return loadbalancer.servers[server_index], nil
}

// Returns the adress of the simple server instance
func (s *SimpleServer) Address() string { return s.addr }

// Serves the through the reverse proxy
func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// Function to calculate the hash value of a given string
func hash(str string) uint32 {
	_hash := fnv.New32a()
	_hash.Write([]byte(str))

	return _hash.Sum32()
}

// Forwards the request to the server returned by the getNextAvailableServer method
func (loadbalancer *Loadbalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	target, err := loadbalancer.getNextAvailableServer(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Forwarding request to adress %q\n", target.Address())

	target.Serve(rw, req)
}
