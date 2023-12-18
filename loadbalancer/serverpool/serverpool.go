package pool

import "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"

type ServerPool interface {
	AddBackEnd(backend backend.Backend)
	GetBackends() int
}
