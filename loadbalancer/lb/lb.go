package lb

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
)

type LbConfig struct {
	Ip          string
	Port        int
	Protocol    string
	Algorithm   string
	BackendUrls []string
}

type LoadBalancer struct {
	algorithm  enums.LoadBalancingAlgorithmType
	ServerPool pool.ServerPool
	//db and cache
}

// Serve: the loadbalancer serves request to the next backend
func (loadbalancer *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	nextServer := loadbalancer.ServerPool.GetNextServer()
	if nextServer != nil {
		nextServer.Serve(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func CreateLB(config LbConfig) (*LoadBalancer, error) {
	var (
		err        error
		serverPool pool.ServerPool
		algorithm  enums.LoadBalancingAlgorithmType
	)

	switch config.Algorithm {
	case "RoundRobin":
		algorithm = enums.RoundRobin
	default:
		log.Fatal("no algorithm configured")
	}

	serverPool, err = configureUrls(algorithm, config.BackendUrls)
	if err != nil {
		return nil, err
	}

	lb := &LoadBalancer{
		algorithm:  algorithm,
		ServerPool: serverPool,
	}

	// server the load balancer on tcp connection
	return lb, nil
}

func configureUrls(algorithm enums.LoadBalancingAlgorithmType, backendUrls []string) (pool.ServerPool, error) {
	serverPool, err := createPool(algorithm)
	if err != nil {
		log.Fatal(err)
	}

	for _, backendUrl := range backendUrls {
		url, err := url.Parse(backendUrl)
		if err != nil {
			log.Fatalln(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := backend.NewBackend(url, proxy)
		backend.SetAlive(true)

		serverPool.AddBackEnd(backend)
	}

	return serverPool, nil
}

func createPool(algorithm enums.LoadBalancingAlgorithmType) (pool.ServerPool, error) {
	switch algorithm {
	case enums.RoundRobin:
		return &pool.RoundRobinPool{}, nil
	default:
		return nil, errors.New("no algorithm configured")
	}
}

func modifyRequest(url url.URL, request *http.Request) error {
	return nil
}
