package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/samber/lo"
)

var _ ServerPool = (*WeightedRoundRobinPool)(nil)

type WeightedRoundRobinPool struct {
	weightedBackends []*WeightedS
}

type WeightedS struct {
	backend         backend.IBackend
	currentWeight   int
	effectiveWeight int
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) AddBackEnd(backend backend.IBackend) {
	ws := &WeightedS{
		backend:         backend,
		effectiveWeight: backend.GetWeight(),
	}
	weightedRoundRobinPool.weightedBackends = append(weightedRoundRobinPool.weightedBackends, ws)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackendCount() int {
	return len(weightedRoundRobinPool.weightedBackends)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackends() []backend.IBackend {
	backends := lo.Map(weightedRoundRobinPool.weightedBackends, func(item *WeightedS, index int) backend.IBackend {
		return item.backend
	})
	return backends
}

// GetNextServer implement smooth weighted round-robin balancing. Check the following link to see it works
// https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1
func (weightedRoundRobinPool *WeightedRoundRobinPool) GetNextServer(ip string) backend.IBackend {
	//increment the respective be with their assigned weight and get the culmulative sum of all weight
	var best *WeightedS
	total := 0

	for i := 0; i < len(weightedRoundRobinPool.weightedBackends); i++ {
		ws := weightedRoundRobinPool.weightedBackends[i]

		if ws == nil || !ws.backend.IsAlive() {
			continue
		}

		ws.currentWeight += ws.effectiveWeight
		total += ws.effectiveWeight
		if ws.effectiveWeight < ws.backend.GetWeight() {
			ws.effectiveWeight++
		}

		if best == nil || ws.currentWeight > best.currentWeight {
			best = ws
		}
	}

	if best == nil {
		return nil
	}

	best.currentWeight -= total
	return best.backend
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = ProxyErrorHandler(proxy, weightedRoundRobinPool, url)

		backend := backend.NewBackend(url, proxy, backend.WithWeight(config.Weight))
		backend.SetAlive(true)
		weightedRoundRobinPool.AddBackEnd(backend)
	}
}
