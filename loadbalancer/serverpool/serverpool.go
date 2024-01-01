package pool

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/helpers"
)

const (
	Attempts int = iota
	Retry
)

type BeConfig struct {
	Url    string
	Weight int
}
type LbConfig struct {
	Ip        string
	Port      int
	Protocol  string
	Algorithm string
	Ring      int
	BeConfigs []BeConfig
}

type ServerPool interface {
	AddBackEnd(backend backend.IBackend)
	GetBackendCount() int
	GetNextServer(ip string) backend.IBackend
	GetBackends() []backend.IBackend
	ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig)
}

func CreatePool(algorithm enums.LoadBalancingAlgorithmType, ringNumber int) (ServerPool, error) {
	switch algorithm {
	case enums.RoundRobin:
		return &RoundRobinPool{}, nil
	case enums.WeightedRoundRobin:
		return &WeightedRoundRobinPool{}, nil
	case enums.LeastConnection:
		return &LeastConnPool{}, nil
	case enums.IpHash:
		return &IPHashPool{
			ring: helpers.NewRing(ringNumber),
		}, nil
	case enums.LeastResponseTime:
		return &LeastResponseTimePool{}, nil
	default:
		return nil, errors.New("no algorithm configured")
	}
}

func ProxyErrorHandler(proxy *httputil.ReverseProxy, pool ServerPool, backendUrl *url.URL) func(w http.ResponseWriter, r *http.Request, e error) {
	return func(w http.ResponseWriter, r *http.Request, e error) {
		fmt.Printf("%s %s\n", r.URL.Host, e.Error())
		retries := GetRetryFromRequestContext(r.Context())
		if retries < 3 {
			// select {
			// case <-time.After(10 * time.Millisecond):
			// 	proxy.ServeHTTP(w, r)
			// }

			time.Sleep(10 * time.Millisecond)
			ctx := context.WithValue(r.Context(), Retry, retries+1)
			proxy.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		//set the backend alive status to be false
		for _, b := range pool.GetBackends() {
			u := b.GetURL()
			if backendUrl.String() == u.String() {
				b.SetAlive(false)
				break
			}
		}

		attempts := GetRetryFromRequestContext(r.Context())
		fmt.Printf("%s(%s) Attempting retry %d\n", r.RemoteAddr, r.URL.Path, attempts)
		ctx := context.WithValue(r.Context(), Attempts, attempts+1)

		pool.GetNextServer(r.WithContext(ctx).Header.Get("X-Client-Ip"))
	}
}

func GetAttemptFromRequestContext(ctx context.Context) int {
	if v, ok := ctx.Value(Attempts).(int); ok {
		return v
	}

	return 1
}

func GetRetryFromRequestContext(ctx context.Context) int {
	if v, ok := ctx.Value(Retry).(int); ok {
		return v
	}

	return 0
}
