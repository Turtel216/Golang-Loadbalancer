package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {

	// Command line flag to specify port number when running go run ./src/main.go
	// Default value is :4000
	port := flag.String("port", "4000", "network port")

	flag.Parse()

	// Target servers
	servers := []Server{
		newSimpleServer("https://www.youtube.com"),
		newSimpleServer("https://www.facebook.com"),
	}

	// Creates a new loadbalancer at port 8000
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
