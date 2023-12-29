package pool

import (
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

var _ ServerPool = (*IPHashPool)(nil)

type IPHashPool struct {
}

func (ipHashPool IPHashPool) AddBackEnd(backend backend.IBackend) {

}

func (ipHashPool IPHashPool) GetBackendCount() int {
	return 0
}

func (ipHashPool IPHashPool) GetNextServer() backend.IBackend {
	return nil
}

func (ipHashPool IPHashPool) GetBackends() []backend.IBackend {
	return nil
}

func (ipHashPool IPHashPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {

}
