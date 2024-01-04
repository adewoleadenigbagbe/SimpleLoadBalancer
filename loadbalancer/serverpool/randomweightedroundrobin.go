package pool

import (
	"log"
	"math/rand"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/samber/lo"
)

var _ ServerPool = (*RandomWeightedRoundRobinPool)(nil)

type RandomWeightedRoundRobinPool struct {
	backends     []backend.IBackend
	sumOfWeights int
	r            *rand.Rand
}

func (randomWeightedRoundRobinPool *RandomWeightedRoundRobinPool) AddBackEnd(backend backend.IBackend) {
	randomWeightedRoundRobinPool.backends = append(randomWeightedRoundRobinPool.backends, backend)
}

func (randomWeightedRoundRobinPool *RandomWeightedRoundRobinPool) GetBackendCount() int {
	return len(randomWeightedRoundRobinPool.backends)
}

func (randomWeightedRoundRobinPool *RandomWeightedRoundRobinPool) GetBackends() []backend.IBackend {
	return randomWeightedRoundRobinPool.backends
}

func (randomWeightedRoundRobinPool *RandomWeightedRoundRobinPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = ProxyErrorHandler(proxy, randomWeightedRoundRobinPool, url)

		backend := backend.NewBackend(url, proxy, backend.WithWeight(config.Weight))
		backend.SetAlive(true)
		randomWeightedRoundRobinPool.AddBackEnd(backend)
		randomWeightedRoundRobinPool.sumOfWeights += backend.GetWeight()
	}
}

func (randomWeightedRoundRobinPool *RandomWeightedRoundRobinPool) GetNextServer(ip string) backend.IBackend {
	healthyBackends := lo.Filter(randomWeightedRoundRobinPool.backends, func(item backend.IBackend, index int) bool {
		return item.IsAlive()
	})

	if len(healthyBackends) == 0 {
		return nil
	}

	sumWeights := lo.SumBy(healthyBackends, func(b backend.IBackend) int {
		return b.GetWeight()
	})

	if sumWeights <= 0 {
		return nil
	}

	randomWeight := randomWeightedRoundRobinPool.r.Intn(sumWeights) + 1
	for _, b := range healthyBackends {
		randomWeight = randomWeight - b.GetWeight()
		if randomWeight <= 0 {
			return b
		}
	}

	return healthyBackends[len(healthyBackends)-1]
}
