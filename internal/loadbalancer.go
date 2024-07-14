package loadbalancer

import "net/http"

//
// Types
//

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}
