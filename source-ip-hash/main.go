package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"log"
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
	addr  string
	proxy *httputil.ReverseProxy
}

// Creates and returns a new instance of the simpleServer struct
func newSimpleServer(addr string) simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	return simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type Loadbalancer struct {
	port    string
	mu      sync.Mutex
	servers []simpleServer
}

// Creates and returns a new loadbalancer instance
func NewLoadBalancer(port string, servers []simpleServer) *Loadbalancer {
	loadbalancer := &Loadbalancer{
		port:    port,
		servers: servers,
	}

	return loadbalancer
}

// Returns the adress of the simple server instance
func (s *simpleServer) Address() string { return s.addr }

// Serves the through the reverse proxy
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// Returns the server selected by the source ip hash algorithm
func (loadbalancer *Loadbalancer) getNextAvailableServer(req *http.Request) (simpleServer, error) {
	// get the source ip
	request_ip := req.Header.Get("X-Forwarded-For")

	// calculate the hash of the source ip
	ip_hash := hash(request_ip)

	// Map the hash value to the server index
	server_index := int(ip_hash) % len(loadbalancer.servers)

	return loadbalancer.servers[server_index], nil
}

// Function to calculate the hash value of a given string
func hash(str string) uint32 {
	_hash := fnv.New32a()
	_hash.Write([]byte(str))

	return _hash.Sum32()
}

// Forwards the request to the server returned by the getNextAvailableServer method
func (loadbalancer *Loadbalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	target, err := loadbalancer.getNextAvailableServer(req)
	if err != nil {
		log.Fatal(err)
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
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.youtube.com"),
		newSimpleServer("https://www.facebook.com"),
	}

	// Create a new loadbalancer
	loadbalancer := NewLoadBalancer(*port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.serveProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.port)
	http.ListenAndServe(":"+loadbalancer.port, nil)
}
