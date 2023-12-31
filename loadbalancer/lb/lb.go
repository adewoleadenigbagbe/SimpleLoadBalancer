package lb

import (
	"errors"
	"net/http"

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
	modifyRequest(r)
	nextServer := loadbalancer.ServerPool.GetNextServer(r.Header.Get("X-Client-Ip"))
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
		err = errors.New("no algorithm configured")
		return nil, err
	}

	serverPool, err = pool.CreatePool(algorithm, ringNumber)
	if err != nil {
		return nil, err
	}

	serverPool.ConfigurePool(algorithm, config.BeConfigs)

	lb := &LoadBalancer{
		algorithm:  algorithm,
		ServerPool: serverPool,
	}
	return lb, nil
}

func modifyRequest(request *http.Request) {
	setHeaders(request)
}

func setHeaders(request *http.Request) {
	request.Header.Set("X-Forwarded", request.Host)
	request.Header.Set("X-Client-Ip", request.RemoteAddr)
}
