package pool

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoundRobinPool(t *testing.T) {
	serverpool, _ := CreatePool(enums.RoundRobin, 0)
	_, ok := serverpool.(*RoundRobinPool)
	assert.Equal(t, true, ok)
}

func TestGetPoolSize(t *testing.T) {
	url, _ := url.Parse("http://localhost:3333")
	serverpool, _ := CreatePool(enums.RoundRobin, 0)
	roundRobinPool, _ := serverpool.(*RoundRobinPool)

	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url)).(*backend.Backend)
	roundRobinPool.AddBackEnd(b)

	//Assert
	assert.Equal(t, 1, roundRobinPool.GetBackendCount())
}
