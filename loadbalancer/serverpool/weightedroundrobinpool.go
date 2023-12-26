package pool

import "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"

var _ ServerPool = (*WeightedRoundRobinPool)(nil)

type WeightedRoundRobinPool struct {
	current  int
	backends []backend.IBackend
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) AddBackEnd(backend backend.IBackend) {
	weightedRoundRobinPool.backends = append(weightedRoundRobinPool.backends, backend)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackendCount() int {
	return len(weightedRoundRobinPool.backends)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetNextServer() backend.IBackend {
	weightedRoundRobinPool.current = (weightedRoundRobinPool.current + 1) % len(weightedRoundRobinPool.backends)
	return weightedRoundRobinPool.backends[weightedRoundRobinPool.current]
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackends() []backend.IBackend {
	return weightedRoundRobinPool.backends
}
