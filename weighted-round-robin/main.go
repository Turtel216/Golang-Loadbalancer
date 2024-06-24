package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

type simpleServer struct {
	addr   string
	weight int
	proxy  *httputil.ReverseProxy
}

// Creates and returns a new instance of the simpleServer struct
func newSimpleServer(addr string, weight int) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	return &simpleServer{
		addr:   addr,
		weight: weight,
		proxy:  httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Loadbalancer struct {
	port            string
	roundRobinCount int
	servers         []Server //Interface
	total_weights   int
}

// Creates and returns a new loadbalancer instance
func NewLoadBalancer(port string, servers []Server, total_weights int) *Loadbalancer {

	return &Loadbalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
		total_weights:   total_weights,
	}
}

// Returns the adress of the simple server instance
func (s *simpleServer) Address() string { return s.addr }

// Ensures that the simpleServer instsance is running
func (s *simpleServer) IsAlive() bool { return true }

// Serves the through the reverse proxy
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// returns the server selected by the weighted round-robin scheduler
func (loadbalancer *Loadbalancer) getNextAvailableServer() (server Server) {
	server = loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]

	for !server.IsAlive() {
		loadbalancer.roundRobinCount++
		server = loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]
	}

	loadbalancer.roundRobinCount++
	return
}

// Forwards the request to the server returned by the getNextAvailableServer method
func (loadbalancer *Loadbalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	target := loadbalancer.getNextAvailableServer()
	fmt.Printf("Forwarding request to adress %q\n", target.Address())

	target.Serve(rw, req)
}

func main() {

	// Command line flag to specify port number when running go run ./src/main.go
	// Default value is :4000
	port := flag.String("port", "4000", "network port")

	flag.Parse()

	// Target servers
	servers := []Server{
		newSimpleServer("https://www.google.com", 3),
		newSimpleServer("https://www.youtube.com", 2),
		newSimpleServer("https://www.facebook.com", 1),
	}

	// Creates a new loadbalancer at port 8000
	loadbalancer := NewLoadBalancer(*port, servers, 6)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.serveProxy(rw, req)
	}

	//Retoutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.port)
	http.ListenAndServe(":"+loadbalancer.port, nil)
}
