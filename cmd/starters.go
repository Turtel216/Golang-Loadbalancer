package main

import (
	"fmt"
	"net/http"

	round_robin "github.com/Turtel216/Golang-Loadbalancer/internal/round-robin"
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

// starts up the round-robin loadbalancer
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
