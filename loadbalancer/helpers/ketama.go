//This piece of code  for ketama consistent hashing was a reference to the link below , it was modified to suit the implementation of
//IP hash for the load balancing algorithm used in this project
// https://github.com/mncaudill/ketama/blob/bea055a2a9ba19da658bf2363a6706288c5f8653/ketama_test.go

package helpers

import (
	"crypto/sha1"
	"sort"
	"strconv"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/loadbalancer/backend"
)

type node struct {
	backend backend.IBackend
	hash    uint
}

type continumPoints []node

func (p continumPoints) Len() int           { return len(p) }
func (p continumPoints) Less(i, j int) bool { return p[i].hash < p[j].hash }
func (p continumPoints) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p continumPoints) Sort()              { sort.Sort(p) }

type HashRing struct {
	defaultSpots int
	points       continumPoints
	length       int
}

func NewRing(n int) (h *HashRing) {
	h = new(HashRing)
	h.defaultSpots = n
	return
}

// Adds a new node to a hash ring
func (h *HashRing) AddNode(b backend.IBackend) {
	weight := b.GetWeight()
	tSpots := h.defaultSpots * weight
	hash := sha1.New()
	for i := 1; i <= tSpots; i++ {
		url := b.GetURL()
		hash.Write([]byte(url.String() + ":" + strconv.Itoa(i)))
		hashBytes := hash.Sum(nil)

		n := &node{
			backend: b,
			hash:    uint(hashBytes[19]) | uint(hashBytes[18])<<8 | uint(hashBytes[17])<<16 | uint(hashBytes[16])<<24,
		}

		h.points = append(h.points, *n)
		hash.Reset()
	}
}

func (h *HashRing) Bake() {
	h.points.Sort()
	h.length = len(h.points)
}

func (h *HashRing) Hash(s string) backend.IBackend {
	hash := sha1.New()
	hash.Write([]byte(s))
	hashBytes := hash.Sum(nil)
	v := uint(hashBytes[19]) | uint(hashBytes[18])<<8 | uint(hashBytes[17])<<16 | uint(hashBytes[16])<<24
	i := sort.Search(h.length, func(i int) bool { return h.points[i].hash >= v })

	if i == h.length {
		i = 0
	}

	return h.points[i].backend
}
