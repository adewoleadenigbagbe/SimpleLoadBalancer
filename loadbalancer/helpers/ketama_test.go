package helpers

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strconv"
	"testing"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
)

func TestGetInfo(t *testing.T) {
	ring := NewRing(200)

	url1, _ := url.Parse("http://localhost:3333")
	n1 := backend.NewBackend(url1, httputil.NewSingleHostReverseProxy(url1), backend.WithWeight(1)).(*backend.Backend)

	url2, _ := url.Parse("http://localhost:3334")
	n2 := backend.NewBackend(url2, httputil.NewSingleHostReverseProxy(url2), backend.WithWeight(1)).(*backend.Backend)

	url3, _ := url.Parse("http://localhost:3335")
	n3 := backend.NewBackend(url3, httputil.NewSingleHostReverseProxy(url3), backend.WithWeight(2)).(*backend.Backend)

	url4, _ := url.Parse("http://localhost:3336")
	n4 := backend.NewBackend(url4, httputil.NewSingleHostReverseProxy(url4), backend.WithWeight(5)).(*backend.Backend)

	nodes := []backend.IBackend{n1, n2, n3, n4}

	for _, node := range nodes {
		ring.AddNode(node)
	}

	ring.Bake()

	m := make(map[backend.IBackend]int)
	for i := 0; i < 1e6; i++ {
		m[ring.Hash("test value"+strconv.FormatUint(uint64(i), 10))]++
	}

	for _, node := range nodes {
		fmt.Println(node, m[node])
	}
}
