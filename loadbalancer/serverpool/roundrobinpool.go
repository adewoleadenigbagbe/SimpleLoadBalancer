package pool

import "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"

var _ ServerPool = (*RoundRobinPool)(nil)

type RoundRobinPool struct {
	backends []backend.Backend
}

func (roundRobinPool RoundRobinPool) AddBackEnd(backend backend.Backend) {
	roundRobinPool.backends = append(roundRobinPool.backends, backend)
}

func (roundRobinPool RoundRobinPool) GetBackends() int {
	return len(roundRobinPool.backends)
}
