package pool

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/helpers"
)

var _ ServerPool = (*IPHashPool)(nil)

type IPHashPool struct {
	backends []backend.IBackend
	ring     *helpers.HashRing
}

func (ipHashPool *IPHashPool) AddBackEnd(backend backend.IBackend) {
	ipHashPool.ring.AddNode(backend)
}

func (ipHashPool *IPHashPool) GetBackendCount() int {
	return len(ipHashPool.backends)
}

func (ipHashPool *IPHashPool) GetNextServer(ip string) backend.IBackend {
	return ipHashPool.ring.Hash(ip)
}

func (ipHashPool *IPHashPool) GetBackends() []backend.IBackend {
	return ipHashPool.backends
}

func (ipHashPool *IPHashPool) ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig) {
	for _, config := range configs {
		url, err := url.Parse(config.Url)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ErrorHandler = ProxyErrorHandler(proxy, ipHashPool, url)

		backend := backend.NewBackend(url, proxy)

		backend.SetAlive(true)
		ipHashPool.AddBackEnd(backend)
	}

	for _, backend := range ipHashPool.backends {
		ipHashPool.ring.AddNode(backend)
	}

	//sort the continum points
	ipHashPool.ring.Bake()
}
