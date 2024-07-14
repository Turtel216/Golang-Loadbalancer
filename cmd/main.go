package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/round-robin"
)

const (
	ROUND_ROBIN       = string("round-robin")
	LEAST_CONNECTIONS = string("least-connections")
)

// Target servers
var servers = []round_robin.Server{
	round_robin.NewSimpleServer("https://www.youtube.com"),
	round_robin.NewSimpleServer("https://www.facebook.com"),
	round_robin.NewSimpleServer("https://www.google.com"),
}

// run the loadbalancer specified by the input string
func start_loadbalancer(algo_type, port *string) error {
	switch *algo_type {
	case ROUND_ROBIN:
		run_round_robin(port)
		return nil
	case LEAST_CONNECTIONS:
		return nil
	default:
		return fmt.Errorf("%s is not a valid algorithm type", *algo_type)
	}
}

// starts up the round-robin loadbalancer
func run_round_robin(port *string) {
	// Creates a new loadbalancer at port 8000
	loadbalancer := round_robin.NewLoadBalancer(*port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}

func main() {

	// Command line flag to specify port number when running go run ./src/main.go
	// Default value is :4000
	port := flag.String("port", "4000", "network port")

	// Command line flag to specify what kind of algorithm to use
	// Default algorithm is round-robin
	algo_type := flag.String("type", "0", "type of load balancing algorithm")

	flag.Parse()

	if err := start_loadbalancer(algo_type, port); err != nil {
		log.Fatal(err)
	}
}
