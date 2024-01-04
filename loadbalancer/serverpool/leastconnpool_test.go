package pool

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/stretchr/testify/assert"
)

func TestLeastConnPoolSize(t *testing.T) {
	url, _ := url.Parse("http://localhost:3333")
	serverpool, _ := CreatePool(enums.LeastConnection, 0)
	leastConnPool, _ := serverpool.(*LeastConnPool)

	b := backend.NewBackend(url, httputil.NewSingleHostReverseProxy(url)).(*backend.Backend)
	b.SetAlive(true)
	leastConnPool.AddBackEnd(b)

	//Assert
	assert.Equal(t, 1, leastConnPool.GetBackendCount())
}
