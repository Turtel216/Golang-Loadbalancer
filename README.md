# Loadbalancer in Golang

### A loadbalancer written in the [Go Programming Language](https://go.dev/) using various algorithms

## Implemented algorithms

- Round-Robin
- Weighted Round-Robin
- Source ip hash
- Least connections
- Least Response Time(currently being worked on)

## Running the server

#### Before running the server make sure that there is a **loadbalancer.config** file in the project directory. Inside of the .config file type out your target server addresses

The server can be started by running the following command

    go run ./cmd 

This will start the server at the default port 4000 using the round robin algorithm. You can add the -port flag to specify the port numebr and the -type flag to specify the algorithm type. You can choose from the following list:

- **rr** for round robin
- **wrr** for weighted round robin
- **lc** for least connections
- **si** for source ip hash
- **lrs** for least response time

i.g.

    go run ./cmd -port=3000 -type=wrr

This will start the server at port 3000 using the weighted round robin algorithm
