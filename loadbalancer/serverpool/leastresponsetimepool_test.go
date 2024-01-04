package pool

import (
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/enums"
	"github.com/stretchr/testify/assert"
)

func TestNextServeLeastActiveConnection(t *testing.T) {
	serverpool, _ := CreatePool(enums.LeastResponseTime, 0)
	leastResponseTimePool, _ := serverpool.(*LeastResponseTimePool)

	url1, _ := url.Parse("http://localhost:3333")
	b1 := backend.NewBackend(url1, httputil.NewSingleHostReverseProxy(url1), backend.WithConnections(1)).(*backend.Backend)
	b1.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b1)

	url2, _ := url.Parse("http://localhost:3334")
	b2 := backend.NewBackend(url2, httputil.NewSingleHostReverseProxy(url2), backend.WithConnections(2)).(*backend.Backend)
	b2.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b2)

	url3, _ := url.Parse("http://localhost:3335")
	b3 := backend.NewBackend(url3, httputil.NewSingleHostReverseProxy(url3), backend.WithConnections(3)).(*backend.Backend)
	b3.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b3)

	next := leastResponseTimePool.GetNextServer("")

	//Assert
	assert.Equal(t, b1, next)
}

func TestNextServeLeastResponseTime(t *testing.T) {
	serverpool, _ := CreatePool(enums.LeastResponseTime, 0)
	leastResponseTimePool, _ := serverpool.(*LeastResponseTimePool)

	url1, _ := url.Parse("http://localhost:3333")
	b1 := backend.NewBackend(url1, httputil.NewSingleHostReverseProxy(url1), backend.WithConnections(1), backend.WithResponseTime(8*time.Second)).(*backend.Backend)
	b1.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b1)

	url2, _ := url.Parse("http://localhost:3334")
	b2 := backend.NewBackend(url2, httputil.NewSingleHostReverseProxy(url2), backend.WithConnections(1), backend.WithResponseTime(4*time.Second)).(*backend.Backend)
	b2.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b2)

	url3, _ := url.Parse("http://localhost:3335")
	b3 := backend.NewBackend(url3, httputil.NewSingleHostReverseProxy(url3), backend.WithConnections(3), backend.WithResponseTime(1*time.Second)).(*backend.Backend)
	b3.SetAlive(true)
	leastResponseTimePool.AddBackEnd(b3)

	next := leastResponseTimePool.GetNextServer("")

	//Assert
	assert.Equal(t, b2, next)
}
