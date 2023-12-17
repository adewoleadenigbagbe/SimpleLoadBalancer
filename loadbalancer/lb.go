package loadbalancer

import "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"

type LoadBalancer struct {
	backends   []Backend
	algorithm  enums.LoadBalancingAlgorithmType
	ServerPool ServerPool
	//db and cache
}

func (loadBalancer LoadBalancer) ConfigureBackend(backend Backend) {
	//add to cache and also a db
}
