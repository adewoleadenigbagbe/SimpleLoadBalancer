package pool

import (
	"errors"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/helpers"
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
