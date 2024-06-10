package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
  Address() string 
  IsAlive() bool 
  Serve(rw http.ResponseWriter,r *http.Request)
}

type simpleServer struct {
  addr string 
  proxy *httputil.ReverseProxy
}

type Loadbalancer struct {
  port string
  roundRobinCount int
  servers []Server
}

func NewLoadBalancer(port string, servers [] Server) *Loadbalancer {
  return &Loadbalancer {
    port: port,
    roundRobinCount: 0,
    servers: servers,
  }
}

func newSimpleServer(addr string) *simpleServer{
  serverUrl, err := url.Parse(addr)
  if err != nil {
    fmt.Printf("error: %v\n", err)
    os.Exit(1)
  }

  return &simpleServer{
    addr: addr,
    proxy: httputil.NewSingleHostReverseProxy(serverUrl),
  }
}

func (server *simpleServer) Adress() string {return server.addr}

func (server *simpleServer) IsAlive() bool { return true}

func (server *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
  server.proxy.ServeHTTP(rw, req)
}

func (loadbalancer *Loadbalancer) getNextAvailableServer() Server {
  server := loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]

  for !server.IsAlive() {
    loadbalancer.roundRobinCount++
    server = loadbalancer.servers[loadbalancer.roundRobinCount%len(loadbalancer.servers)]
  }

  loadbalancer.roundRobinCount++
  return server
}

func (loadbalancer *Loadbalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
  target := loadbalancer.getNextAvailableServer()
  fmt.Printf("Forwarding request to adress %q\n", target.Adress)
  
  target.Serve(rw, req)
}

func main() {
  servers := []Server {
    newSimpleServer("https://www.google.com"),
    newSimpleServer("https://www.youtube.com"),
    newSimpleServer("https://www.go.dev"),
  }

  loadbalancer := NewLoadBalancer("8000", servers)

  handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
    loadbalancer.serveProxy(rw, req)
  }
  http.HandleFunc("/", handleRedirect)

  fmt.Printf("Loadbalancer started at port %s \n", loadbalancer.port)
  http.ListenAndServe(":" + loadbalancer.port, nil)
}
