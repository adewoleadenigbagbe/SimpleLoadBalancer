package pool

import "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"

var _ ServerPool = (*RoundRobinPool)(nil)

type RoundRobinPool struct {
	current  int
	backends []backend.IBackend
}

func (roundRobinPool *RoundRobinPool) AddBackEnd(backend backend.IBackend) {
	roundRobinPool.backends = append(roundRobinPool.backends, backend)
}

func (roundRobinPool *RoundRobinPool) GetBackends() int {
	return len(roundRobinPool.backends)
}

func (roundRobinPool *RoundRobinPool) GetNextServer() backend.IBackend {
	roundRobinPool.current = (roundRobinPool.current + 1) % len(roundRobinPool.backends)
	return roundRobinPool.backends[roundRobinPool.current]
}
