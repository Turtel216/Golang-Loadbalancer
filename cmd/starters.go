package main

import (
	"fmt"
	"net/http"

	round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/round-robin"
	source_ip_hash "github.com/Turtel216/Golang-Loadbalancer/internal/source-ip-hash"
	weighted_round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/weighted-round-robin"
)

// starts up the round-robin loadbalancer
func run_round_robin(port *string) {
	// Target servers
	servers := []round_robin.Server{
		round_robin.NewSimpleServer("https://www.youtube.com"),
		round_robin.NewSimpleServer("https://www.facebook.com"),
		round_robin.NewSimpleServer("https://www.google.com"),
	}

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

// starts up the weighted round-robin loadbalancer
func run_weighted_round_robin(port *string) {
	// Target servers
	servers := []weighted_round_robin.SimpleServer{
		weighted_round_robin.NewSimpleServer("https://www.youtube.com", 1),
		weighted_round_robin.NewSimpleServer("https://www.facebook.com", 2),
		weighted_round_robin.NewSimpleServer("https://www.google.com", 3),
	}

	// Creates a new loadbalancer at port 8000
	loadbalancer := weighted_round_robin.NewLoadBalancer(*port, servers)

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
func run_source_ip_hash(port *string) {
	// Target servers
	servers := []source_ip_hash.SimpleServer{
		source_ip_hash.NewSimpleServer("https://www.youtube.com"),
		source_ip_hash.NewSimpleServer("https://www.facebook.com"),
		source_ip_hash.NewSimpleServer("https://www.google.com"),
	}

	// Create a new loadbalancer
	loadbalancer := source_ip_hash.NewLoadBalancer(*port, servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadbalancer.ServeProxy(rw, req)
	}

	//Reroutes request coming in at the `/` endpoint
	http.HandleFunc("/", handleRedirect)

	// Starting the server
	fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.Port)
	http.ListenAndServe(":"+loadbalancer.Port, nil)
}
