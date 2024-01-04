package lb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type LoadBalancer struct {
	algorithm  enums.LoadBalancingAlgorithmType
	ServerPool pool.ServerPool
	//db and cache
}

// Serve: the loadbalancer serves request to the next backend
func (loadbalancer *LoadBalancer) Serve(w http.ResponseWriter, r *http.Request) {
	loadbalancer.modifyRequest(r)
	nextServer := loadbalancer.ServerPool.GetNextServer(r.Header.Get("X-Client-Ip"))
	if nextServer != nil {
		nextServer.Serve(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func (loadbalancer *LoadBalancer) HealthCheck(ctx context.Context) {
	healthCheckTicker := time.NewTicker(2 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			healthCheckTicker.Stop()
			fmt.Println("Gracefully shutting down health check")
			return
		case t := <-healthCheckTicker.C:
			fmt.Println("Tick at", t)
			healthyBackends := lo.Filter(loadbalancer.ServerPool.GetBackends(), func(item backend.IBackend, index int) bool {
				return item.IsAlive()
			})
			check(ctx, healthyBackends)
		}
	}
}

func (loadbalancer *LoadBalancer) modifyRequest(request *http.Request) {
	setHeaders(request)
}

func check(ctx context.Context, backends []backend.IBackend) {
	//TODO:might need to change this to zero, but will test first
	aliveChan := make(chan bool, 1)
	for _, b := range backends {
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()
		status := "up"
		backend.IsBackendAlive(requestCtx, aliveChan, b.GetURL())
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully shutting down health check")
			return
		case alive := <-aliveChan:
			b.SetAlive(alive)
			if !alive {
				status = "down"
			}
		}

		url := b.GetURL()
		fmt.Println("URL Status ", url.String(), "is ", status)
	}
}

func setHeaders(request *http.Request) {
	request.Header.Set("X-Request-Id", uuid.New().String())
	request.Header.Set("X-Forwarded", request.Host)
	request.Header.Set("X-Client-Ip", request.RemoteAddr)
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
	case "SmoothWeightedRoundRobin":
		algorithm = enums.SmoothWeightedRoundRobin
	case "RandomWeightedRoundRobin":
		algorithm = enums.RandomWeightedRoundRobin
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
