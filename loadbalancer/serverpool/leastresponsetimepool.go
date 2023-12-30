package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

var _ ServerPool = (*LeastConnPool)(nil)

type LeastResponseTimePool struct {
	backends []backend.IBackend
}

func (leastResponseTimePool *LeastResponseTimePool) AddBackEnd(backend backend.IBackend) {
	leastResponseTimePool.backends = append(leastResponseTimePool.backends, backend)
}

func (leastResponseTimePool *LeastResponseTimePool) GetBackendCount() int {
	return len(leastResponseTimePool.backends)
}

func (leastResponseTimePool *LeastResponseTimePool) GetNextServer(ip string) backend.IBackend {
	var leastResponseTimePeer backend.IBackend
	for _, b := range leastResponseTimePool.backends {
		if b.IsAlive() {
			leastResponseTimePeer = b
			break
		}
	}

	for _, b := range leastResponseTimePool.backends {
		if !b.IsAlive() {
			continue
		}

		if leastResponseTimePeer.GetActiveConnections() >= b.GetActiveConnections() {
			if leastResponseTimePeer.GetActiveConnections() == b.GetActiveConnections() {
				if leastResponseTimePeer.GetResponseTime() > b.GetResponseTime() {
					leastResponseTimePeer = b
					continue
				}
			}
			leastResponseTimePeer = b
		}
	}
	return leastResponseTimePeer
}

func (leastResponseTimePool *LeastResponseTimePool) GetBackends() []backend.IBackend {
	return leastResponseTimePool.backends
}

func (leastResponseTimePool *LeastResponseTimePool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := backend.NewBackend(url, proxy, backend.WithConnections(0))

		backend.SetAlive(true)
		leastResponseTimePool.AddBackEnd(backend)
	}
}
