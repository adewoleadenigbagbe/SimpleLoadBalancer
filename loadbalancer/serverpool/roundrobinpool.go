package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/samber/lo"
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

func (roundRobinPool *RoundRobinPool) GetNextServer(ip string) backend.IBackend {
	healthyBackends := lo.Filter(roundRobinPool.backends, func(item backend.IBackend, index int) bool {
		return item.IsAlive()
	})
	roundRobinPool.current = (roundRobinPool.current + 1) % len(healthyBackends)
	return healthyBackends[roundRobinPool.current]
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
		proxy.ErrorHandler = ProxyErrorHandler(proxy, roundRobinPool, url)

		backend := backend.NewBackend(url, proxy)
		backend.SetAlive(true)
		roundRobinPool.AddBackEnd(backend)
	}
}
