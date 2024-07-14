# Loadbalancer in Golang

### A loadbalancer written in the [Go Programming Language](https://go.dev/) using various algorithms

## Implemented algorithms

- Round-Robin
- Weighted Round-Robin
- Source ip hash
- Least connections
- Least Response Time(currently being worked on)
- **More to be added soon!**

## Running the server

The server can be started by running the following command

    go run ./cmd/main.go 

This will start the server at the default port 4000 using the round robin algorithm. You can add the -port flag to specify the port numebr and the -type flag to specify the algorithm type. You can choose from the following list:

- **rr** for round robin
- **wrr** for weighted round robin
- **lc** for least connections
- **si** for source ip hash
- **lrs** for least response time

i.g.

    go run ./cmd/main.go -port=3000 -type=wrr

This will start the server at port 3000 using the weighted round robin algorithm
