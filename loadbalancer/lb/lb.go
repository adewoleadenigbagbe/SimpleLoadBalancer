package lb

import (
	"log"
	"net/http"
	"net/url"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
)

type LoadBalancer struct {
	algorithm  enums.LoadBalancingAlgorithmType
	ServerPool pool.ServerPool
	//db and cache
}

// Serve: the loadbalancer serves request to the next backend
func (loadbalancer *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	nextServer := loadbalancer.ServerPool.GetNextServer(r.RemoteAddr)
	if nextServer != nil {
		nextServer.Serve(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func CreateLB(config pool.LbConfig) (*LoadBalancer, error) {
	var (
		err        error
		serverPool pool.ServerPool
		algorithm  enums.LoadBalancingAlgorithmType
		ringNumber int
	)

	switch config.Algorithm {
	case "RoundRobin":
		algorithm = enums.RoundRobin
	case "WeightedRoundRobin":
		algorithm = enums.WeightedRoundRobin
	case "IPHash":
		algorithm = enums.IpHash
	case "LeastConnection":
		algorithm = enums.LeastConnection
	case "LeastResponseTime":
		algorithm = enums.LeastResponseTime
		ringNumber = config.Ring
	case "ResourceLoad":
		algorithm = enums.ResourceLoad
	default:
		log.Fatal("no algorithm configured")
	}

	serverPool, err = pool.CreatePool(algorithm, ringNumber)
	if err != nil {
		log.Fatal(err)
	}

	serverPool.ConfigurePool(algorithm, config.BeConfigs)

	lb := &LoadBalancer{
		algorithm:  algorithm,
		ServerPool: serverPool,
	}

	// server the load balancer on tcp connection
	return lb, nil
}

func modifyRequest(url url.URL, request *http.Request) error {
	return nil
}
