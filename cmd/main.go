package main

import (
	"flag"
	"fmt"
	"log"
)

const (
	ROUND_ROBIN       = string("round robin")
	LEAST_CONNECTIONS = string("least connections")
	SOURCE_IP_HASH    = string("source ip")
)

// run the loadbalancer specified by the input string
func start_loadbalancer(algo_type, port *string) error {
	switch *algo_type {
	case ROUND_ROBIN:
		run_round_robin(port)
		return nil
	case LEAST_CONNECTIONS:
		run_weighted_round_robin(port)
		return nil
	case SOURCE_IP_HASH:
		run_source_ip_hash(port)
		return nil
	default:
		return fmt.Errorf("%s is not a valid algorithm type", *algo_type)
	}
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
