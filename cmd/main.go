package main

import (
	"flag"
	"fmt"
	"net/http"

	round_robin "github.com/Turtel216/Golang-Loadbalancer/round-robin"
)

func main() {

	// Command line flag to specify port number when running go run ./src/main.go
	// Default value is :4000
	port := flag.String("port", "4000", "network port")

	flag.Parse()

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
