package pool

import (
	"errors"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
)

type ServerPool interface {
	AddBackEnd(backend backend.IBackend)
	GetBackendCount() int
	GetNextServer() backend.IBackend
	GetBackends() []backend.IBackend
}

func CreatePool(algorithm enums.LoadBalancingAlgorithmType) (ServerPool, error) {
	switch algorithm {
	case enums.RoundRobin:
		return &RoundRobinPool{}, nil
	default:
		return nil, errors.New("no algorithm configured")
	}
}
