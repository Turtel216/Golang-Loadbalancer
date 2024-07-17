package main

import (
	"fmt"
	"net/http"

	lb "github.com/Turtel216/Golang-Loadbalancer/internal"
	least_connections "github.com/Turtel216/Golang-Loadbalancer/internal/least-connections"
	least_response_time "github.com/Turtel216/Golang-Loadbalancer/internal/least-response-time"
	round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/round-robin"
	source_ip_hash "github.com/Turtel216/Golang-Loadbalancer/internal/source-ip-hash"
	weighted_round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/weighted-round-robin"
)

// starts up the round-robin loadbalancer
func run_round_robin(port string, urls []string) {
	//Initialize target servers
	var servers []lb.Server
	for _, url := range urls {
		servers = append(servers, round_robin.NewSimpleServer(url))
	}

	// Creates a new loadbalancer at port 8000
	loadbalancer := round_robin.NewLoadBalancer(port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}

// starts up the weighted round-robin loadbalancer
func run_weighted_round_robin(port string, urls []string) {
	// Target servers
	servers := []weighted_round_robin.SimpleServer{
		*weighted_round_robin.NewSimpleServer("https://www.youtube.com", 1),
		*weighted_round_robin.NewSimpleServer("https://www.facebook.com", 2),
		*weighted_round_robin.NewSimpleServer("https://www.google.com", 3),
	}

	// Creates a new loadbalancer at port 8000
	loadbalancer := weighted_round_robin.NewLoadBalancer(port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}

// starts up the source ip hash loadbalancer
func run_source_ip_hash(port string, urls []string) {
	//Initialize target servers
	var servers []source_ip_hash.SimpleServer
	for _, url := range urls {
		servers = append(servers, *source_ip_hash.NewSimpleServer(url))
	}

	// Create a new loadbalancer
	loadbalancer := source_ip_hash.NewLoadBalancer(port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}

// starts up the least connections loadbalancer
func run_least_connections(port string, urls []string) {
	//Initialize target servers
	var servers []lb.Server
	for _, url := range urls {
		servers = append(servers, least_connections.NewSimpleServer(url))
	}

	// Create a new loadbalancer
	loadbalancer := least_connections.NewLoadBalancer(port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}

// starts up the least connections loadbalancer
func run_least_response_time(port string, urls []string) {
	//Initialize target servers
	var servers []lb.Server
	for _, url := range urls {
		servers = append(servers, least_response_time.NewSimpleServer(url))
	}

	// Create a new loadbalancer
	loadbalancer := least_response_time.NewLoadBalancer(port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}
