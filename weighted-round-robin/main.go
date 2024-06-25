package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
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
func newSimpleServer(addr string, weight int) simpleServer {
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
	port            string
	roundRobinCount int
	servers         []simpleServer
	current         int
	total_weights   int
}

// Creates and returns a new loadbalancer instance
func NewLoadBalancer(port string, servers []simpleServer, total_weights int) *Loadbalancer {
	loadbalancer := &Loadbalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
		total_weights:   total_weights,
	}

	for _, server := range servers {
		loadbalancer.total_weights += server.weight
	}

	return loadbalancer
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
func (loadbalancer *Loadbalancer) getNextAvailableServer() (simpleServer, error) {
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
func (loadbalancer *Loadbalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	target, err := loadbalancer.getNextAvailableServer()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Forwarding request to adress %q\n", target.Address())

	target.Serve(rw, req)
}

func main() {

	// Command line flag to specify port number when running go run ./src/main.go
	// Default value is :4000
	port := flag.String("port", "4000", "network port")

	flag.Parse()

	// Target servers
	servers := []simpleServer{
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
