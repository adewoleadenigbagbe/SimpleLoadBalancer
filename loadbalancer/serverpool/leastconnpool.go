package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

var _ ServerPool = (*LeastConnPool)(nil)

type LeastConnPool struct {
	backends []backend.IBackend
}

func (leastConnPool *LeastConnPool) AddBackEnd(backend backend.IBackend) {
	leastConnPool.backends = append(leastConnPool.backends, backend)
}

func (leastConnPool *LeastConnPool) GetBackendCount() int {
	return len(leastConnPool.backends)
}

func (leastConnPool *LeastConnPool) GetNextServer() backend.IBackend {
	var leastConnectedPeer backend.IBackend
	for _, b := range leastConnPool.backends {
		if b.IsAlive() {
			leastConnectedPeer = b
			break
		}
	}

	for _, b := range leastConnPool.backends {
		if !b.IsAlive() {
			continue
		}
		if leastConnectedPeer.GetActiveConnections() > b.GetActiveConnections() {
			leastConnectedPeer = b
		}
	}
	return leastConnectedPeer
}

func (leastConnPool *LeastConnPool) GetBackends() []backend.IBackend {
	return leastConnPool.backends
}

func (leastConnPool *LeastConnPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := backend.NewBackend(url, proxy, backend.WithConnections(0))

		backend.SetAlive(true)
		leastConnPool.AddBackEnd(backend)
	}
}
