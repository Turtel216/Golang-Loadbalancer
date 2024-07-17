package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Turtel216/Golang-Loadbalancer/internal/util"
)

const (
	ROUND_ROBIN          = "rr"
	WEIGHTED_ROUND_ROBIN = "wrr"
	LEAST_CONNECTIONS    = "lc"
	SOURCE_IP_HASH       = "si"
	LEAST_RESPONSE_TIME  = "lrs"
)

// run the loadbalancer specified by the input string
func start_loadbalancer(algo_type, port *string, urls *[]string) error {
	switch *algo_type {
	case ROUND_ROBIN:
		run_round_robin(port, urls)
		return nil
	case WEIGHTED_ROUND_ROBIN:
		run_weighted_round_robin(port, urls)
		return nil
	case LEAST_CONNECTIONS:
		run_weighted_round_robin(port, urls)
		return nil
	case SOURCE_IP_HASH:
		run_source_ip_hash(port, urls)
		return nil
	case LEAST_RESPONSE_TIME:
		run_least_response_time(port, urls)
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
	algo_type := flag.String("type", "rr", "type of load balancing algorithm")

	flag.Parse()

	var path string = "loadbalancer.config"
	urls, err := util.Config_parser(&path)
	if err != nil {
		log.Fatal(err)
	}

	if err = start_loadbalancer(algo_type, port, urls); err != nil {
		log.Fatal(err)
	}
}
