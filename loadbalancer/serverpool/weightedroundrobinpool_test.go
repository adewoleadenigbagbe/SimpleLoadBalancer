package pool

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
	"github.com/stretchr/testify/assert"
)

func TestGetNextServer(t *testing.T) {
	url1, _ := url.Parse("http://localhost:3330")
	b1 := backend.NewBackend(url1, httputil.NewSingleHostReverseProxy(url1), backend.WithWeight(4)).(*backend.Backend)

	ws1 := &WeightedS{
		backend:         b1,
		currentWeight:   0,
		effectiveWeight: b1.GetWeight(),
	}

	url2, _ := url.Parse("http://localhost:3331")
	b2 := backend.NewBackend(url2, httputil.NewSingleHostReverseProxy(url2), backend.WithWeight(3)).(*backend.Backend)

	ws2 := &WeightedS{
		backend:         b2,
		currentWeight:   0,
		effectiveWeight: b2.GetWeight(),
	}

	url3, _ := url.Parse("http://localhost:3332")
	b3 := backend.NewBackend(url3, httputil.NewSingleHostReverseProxy(url3), backend.WithWeight(2)).(*backend.Backend)

	ws3 := &WeightedS{
		backend:         b3,
		currentWeight:   0,
		effectiveWeight: b3.GetWeight(),
	}

	url4, _ := url.Parse("http://localhost:3333")
	b4 := backend.NewBackend(url4, httputil.NewSingleHostReverseProxy(url4), backend.WithWeight(1)).(*backend.Backend)

	ws4 := &WeightedS{
		backend:         b4,
		currentWeight:   0,
		effectiveWeight: b4.GetWeight(),
	}

	wrr := &WeightedRoundRobinPool{
		weightedBackends: []*WeightedS{
			ws1,
			ws2,
			ws3,
			ws4,
		},
	}

	results := make(map[backend.IBackend]int)
	noofTimes := 1000
	for i := 0; i < noofTimes; i++ {
		b := wrr.GetNextServer()
		results[b]++
	}

	totalWeight := b1.GetWeight() + b2.GetWeight() + b3.GetWeight() + b4.GetWeight()
	expectedCount1 := (b1.GetWeight() * noofTimes) / totalWeight
	expectedCount2 := (b2.GetWeight() * noofTimes) / totalWeight
	expectedCount3 := (b3.GetWeight() * noofTimes) / totalWeight
	expectedCount4 := (b4.GetWeight() * noofTimes) / totalWeight

	//Assert
	assert.Equal(t, expectedCount1, results[b1])
	assert.Equal(t, expectedCount2, results[b2])
	assert.Equal(t, expectedCount3, results[b3])
	assert.Equal(t, expectedCount4, results[b4])
}
