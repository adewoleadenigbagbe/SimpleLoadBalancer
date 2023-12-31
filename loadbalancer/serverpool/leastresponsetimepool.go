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
	var current backend.IBackend
	for _, b := range leastResponseTimePool.backends {
		if b.IsAlive() {
			current = b
			break
		}
	}

	for _, b := range leastResponseTimePool.backends {
		if !b.IsAlive() {
			continue
		}

		leastActiveConn := current.GetActiveConnections()
		bActiveConn := b.GetActiveConnections()
		if leastActiveConn >= bActiveConn {
			if leastActiveConn == bActiveConn {
				avgA := float64(current.GetResponseTime()) / float64(leastActiveConn)
				avgB := float64(b.GetResponseTime()) / float64(bActiveConn)
				if avgA > avgB {
					current = b
					continue
				}
			}
			current = b
		}
	}
	return current
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
