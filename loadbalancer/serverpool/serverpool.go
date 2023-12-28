package pool

import (
	"errors"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
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
	BeConfigs []BeConfig
}

type ServerPool interface {
	AddBackEnd(backend backend.IBackend)
	GetBackendCount() int
	GetNextServer() backend.IBackend
	GetBackends() []backend.IBackend
	ConfigurePool(algorithm enums.LoadBalancingAlgorithmType, configs []BeConfig)
}

func CreatePool(algorithm enums.LoadBalancingAlgorithmType) (ServerPool, error) {
	switch algorithm {
	case enums.RoundRobin:
		return &RoundRobinPool{}, nil
	case enums.WeightedRoundRobin:
		return &WeightedRoundRobinPool{}, nil
	default:
		return nil, errors.New("no algorithm configured")
	}
}
