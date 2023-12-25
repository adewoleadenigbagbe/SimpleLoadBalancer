package lb

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoundRobinPool(t *testing.T) {
	serverpool, _ := createPool(enums.RoundRobin)
	_, ok := serverpool.(*pool.RoundRobinPool)
	assert.Equal(t, true, ok)
}

func TestGetPoolSize(t *testing.T) {
	serverpool, _ := createPool(enums.RoundRobin)
	roundRobinPool, _ := serverpool.(*pool.RoundRobinPool)
	url, _ := url.Parse("http://localhost:3333")
	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url)).(*backend.Backend)
	roundRobinPool.AddBackEnd(b)

	assert.Equal(t, 1, roundRobinPool.GetBackends())
}

// func TestModifyRequest(t *testing.T) {

// }

// func TestServe(t *testing.T) {

// }
