package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct {
  addr string 
  proxy *httputil.ReverseProxy
}

type Server interface {
  Address() string 
  IsAlive() bool 
  Serve(rw http.ResponseWriter,r *http.Request)
}

type Loadbalancer struct {
  port string
  roundRobinCount int
  servers []Server
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
