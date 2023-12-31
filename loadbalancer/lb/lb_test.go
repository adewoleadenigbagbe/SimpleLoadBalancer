package lb

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	pool "github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/serverpool"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func RequestResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello there"))
}

var (
	h = http.HandlerFunc(RequestResponseHandler)
)

func TestNotServeWhenAlgorithmIsNotConfigured(t *testing.T) {
	//Mock Servers
	noofServers := 5
	var urls []*url.URL
	for i := 0; i < noofServers; i++ {
		mockServer := httptest.NewServer(h)
		defer mockServer.Close()
		serverUrl, err := url.Parse(mockServer.URL)
		urls = append(urls, serverUrl)
		if err != nil {
			t.Fatal(err)
		}
	}

	urlstrs := lo.Map(urls, func(item *url.URL, index int) string {
		return item.String()
	})

	//Load Balancer
	var beConfigs []pool.BeConfig
	for _, urlstr := range urlstrs {
		beConfig := pool.BeConfig{
			Url:    urlstr,
			Weight: 0,
		}
		beConfigs = append(beConfigs, beConfig)
	}

	config := pool.LbConfig{
		Algorithm: "None",
		Ip:        "localhost",
		Port:      3662,
		Protocol:  "http",
		BeConfigs: beConfigs,
	}

	_, err := CreateLB(config)
	expected := errors.New("no algorithm configured")
	assert.Equal(t, expected, err)

}

func TestServeRobinPool(t *testing.T) {
	//Mock Servers
	noofServers := 5
	var urls []*url.URL
	for i := 0; i < noofServers; i++ {
		mockServer := httptest.NewServer(h)
		defer mockServer.Close()
		serverUrl, err := url.Parse(mockServer.URL)
		urls = append(urls, serverUrl)
		if err != nil {
			t.Fatal(err)
		}
	}

	urlstrs := lo.Map(urls, func(item *url.URL, index int) string {
		return item.String()
	})

	//Load Balancer
	var beConfigs []pool.BeConfig
	for _, urlstr := range urlstrs {
		beConfig := pool.BeConfig{
			Url:    urlstr,
			Weight: 0,
		}
		beConfigs = append(beConfigs, beConfig)
	}

	config := pool.LbConfig{
		Algorithm: "RoundRobin",
		Ip:        "localhost",
		Port:      3662,
		Protocol:  "http",
		BeConfigs: beConfigs,
	}

	lb, _ := CreateLB(config)
	lbServer := httptest.NewServer(http.HandlerFunc(lb.Serve))
	defer lbServer.Close()

	//Client
	noOfRequest := 1000
	loadBalancerUrl, _ := url.Parse(lbServer.URL)
	for i := 0; i < noOfRequest; i++ {
		req, _ := http.NewRequest("GET", loadBalancerUrl.String(), nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		io.ReadAll(res.Body)
	}

	var connections []int
	backends := lb.ServerPool.GetBackends()
	for _, b := range backends {
		connections = append(connections, b.GetActiveConnections())
	}

	t.Log(connections)

	avgConnections := noOfRequest / noofServers
	b := lo.EveryBy(connections, func(x int) bool {
		return x == avgConnections
	})

	assert.Equal(t, true, b)
}

func TestServeWeightedRobinPool(t *testing.T) {
	//Mock Servers
	noofServers := 5
	var urls []*url.URL
	for i := 0; i < noofServers; i++ {
		mockServer := httptest.NewServer(h)
		defer mockServer.Close()
		serverUrl, err := url.Parse(mockServer.URL)
		urls = append(urls, serverUrl)
		if err != nil {
			t.Fatal(err)
		}
	}

	urlstrs := lo.Map(urls, func(item *url.URL, index int) string {
		return item.String()
	})

	//Load Balancer
	var beConfigs []pool.BeConfig
	initialweight := 0
	for _, urlstr := range urlstrs {
		initialweight++
		beConfig := pool.BeConfig{
			Url:    urlstr,
			Weight: initialweight,
		}
		beConfigs = append(beConfigs, beConfig)
	}

	config := pool.LbConfig{
		Algorithm: "WeightedRoundRobin",
		Ip:        "localhost",
		Port:      3662,
		Protocol:  "http",
		BeConfigs: beConfigs,
	}

	lb, _ := CreateLB(config)
	lbServer := httptest.NewServer(http.HandlerFunc(lb.Serve))
	defer lbServer.Close()

	//Client
	noOfRequest := 1000
	loadBalancerUrl, _ := url.Parse(lbServer.URL)
	for i := 0; i < noOfRequest; i++ {
		req, _ := http.NewRequest("GET", loadBalancerUrl.String(), nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		io.ReadAll(res.Body)
	}

	totalWeight := 0
	for _, b := range lb.ServerPool.GetBackends() {
		totalWeight += b.GetWeight()
	}

	t.Log("total weight :", totalWeight)
	for _, b := range lb.ServerPool.GetBackends() {
		v := float64(b.GetWeight()*noOfRequest) / float64(totalWeight)
		connections := math.Round(v)
		assert.Equal(t, connections, float64(b.GetActiveConnections()))
	}
}

func TestServerLeastConnPool(t *testing.T) {

}
