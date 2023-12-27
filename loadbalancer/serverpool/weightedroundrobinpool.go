package pool

import (
	"log"
	"math/rand"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

var _ ServerPool = (*WeightedRoundRobinPool)(nil)

type WeightedRoundRobinPool struct {
	backends         []backend.IBackend
	backendWeightMap []map[string]int
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) AddBackEnd(backend backend.IBackend) {
	weightedRoundRobinPool.backends = append(weightedRoundRobinPool.backends, backend)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackendCount() int {
	return len(weightedRoundRobinPool.backends)
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) GetBackends() []backend.IBackend {
	return weightedRoundRobinPool.backends
}

// GetNextServer implement smooth weighted round-robin balancing. Check the following link to see it works
// https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1
func (weightedRoundRobinPool *WeightedRoundRobinPool) GetNextServer() backend.IBackend {
	//increment the respectiver be with their assigned weight and get the culmulative sum of all weight
	culmulativeWeight := 0
	for i, backend := range weightedRoundRobinPool.backends {
		id := backend.GetID()
		weight := backend.GetWeight()
		weightedRoundRobinPool.backendWeightMap[i][id] += weight
		culmulativeWeight += weightedRoundRobinPool.backendWeightMap[i][id]
	}

	//randomly pick a backend
	index := rand.Intn(len(weightedRoundRobinPool.backends))
	b := weightedRoundRobinPool.backends[index]
	id := b.GetID()
	weightedRoundRobinPool.backendWeightMap[index][id] += culmulativeWeight

	return b
}

func (weightedRoundRobinPool *WeightedRoundRobinPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)

		backend := backend.NewBackend(url, proxy, backend.WithWeight(config.Weight))
		backend.SetAlive(true)
		weightedRoundRobinPool.AddBackEnd(backend)
		id := backend.GetID()
		mp := map[string]int{
			id: 0,
		}
		weightedRoundRobinPool.backendWeightMap = append(weightedRoundRobinPool.backendWeightMap, mp)
	}
}
