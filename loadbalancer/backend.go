package loadbalancer

import (
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

type LoadBalancer struct {
	backends  []Backend
	algorithm enums.LoadBalancingAlgorithmType
}

func (loadBalancer LoadBalancer) ConfigureBackend(backend Backend) {
	//add to cache and also a db
}

type IBackend interface {
	SetAlive(bool)
	IsAlive() bool
	GetURL() *url.URL
	//GetActiveConnections() int
	//Serve(http.ResponseWriter, *http.Request)
}

type Backend struct {
	id       string
	ip       string
	port     int
	protocol string
	alive    bool
}
