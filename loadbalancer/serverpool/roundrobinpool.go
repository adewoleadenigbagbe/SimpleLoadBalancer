package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

var _ ServerPool = (*RoundRobinPool)(nil)

type RoundRobinPool struct {
	current  int
	backends []backend.IBackend
}

func (roundRobinPool *RoundRobinPool) AddBackEnd(backend backend.IBackend) {
	roundRobinPool.backends = append(roundRobinPool.backends, backend)
}

func (roundRobinPool *RoundRobinPool) GetBackendCount() int {
	return len(roundRobinPool.backends)
}

func (roundRobinPool *RoundRobinPool) GetNextServer() backend.IBackend {
	roundRobinPool.current = (roundRobinPool.current + 1) % len(roundRobinPool.backends)
	return roundRobinPool.backends[roundRobinPool.current]
}

func (roundRobinPool *RoundRobinPool) GetBackends() []backend.IBackend {
	return roundRobinPool.backends
}

func (roundRobinPool *RoundRobinPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := backend.NewBackend(url, proxy)

		backend.SetAlive(true)
		roundRobinPool.AddBackEnd(backend)
	}
}
