# Loadbalancer in Golang

### A loadbalancer written in the [Go Programming Language](https://go.dev/) using various algorithms

## Implemented algorithms

- Round-Robin
- Weighted Round-Robin
- Source ip hash (currently being worked on)
- **More to be added soon!**

## Running the server

The server can be started by running the following command

    go run ./src/main.go 

This will start the server at the default port 4000. You can add the -port flag to specify the port numebr

i.g.

    go run ./src/main.go -port=3000

This will start the server at port 3000
