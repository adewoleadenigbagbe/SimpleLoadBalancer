package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/google/uuid"
)

var _ IBackend = (*Backend)(nil)

type IBackend interface {
	GetID() string
	GetWeight() int
	SetAlive(bool)
	IsAlive() bool
	GetURL() url.URL
	GetActiveConnections() int
	Serve(w http.ResponseWriter, r *http.Request)
}

type Metrics struct {
	connections int
	weight      int
}

type MetricsOption func(*Metrics)

type Backend struct {
	id           string
	url          url.URL
	mux          sync.RWMutex
	alive        bool
	metrics      *Metrics
	reverseProxy *httputil.ReverseProxy
}

func (backend *Backend) GetID() string {
	return backend.id
}

func (backend *Backend) SetAlive(alive bool) {
	backend.mux.Lock()
	backend.alive = alive
	backend.mux.Unlock()
}

func (backend *Backend) GetURL() url.URL {
	return backend.url
}

func (backend *Backend) IsAlive() bool {
	backend.mux.RLock()
	defer backend.mux.RUnlock()
	return backend.alive
}

func (backend *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	backend.mux.Lock()
	defer backend.mux.Unlock()
	backend.metrics.connections++
	backend.reverseProxy.ServeHTTP(w, r)
}

func (backend *Backend) GetActiveConnections() int {
	backend.mux.RLock()
	defer backend.mux.RUnlock()
	return backend.metrics.connections
}

func (backend *Backend) GetWeight() int {
	backend.mux.RLock()
	defer backend.mux.RUnlock()
	return backend.metrics.weight
}

func NewBackend(endpoint *url.URL, proxy *httputil.ReverseProxy, options ...MetricsOption) IBackend {
	backend := Backend{
		id:           uuid.NewString(),
		url:          *endpoint,
		reverseProxy: proxy,
		metrics:      &Metrics{},
	}

	for _, opt := range options {
		opt(backend.metrics)
	}
	return &backend
}

func WithConnections(connections int) MetricsOption {
	return func(m *Metrics) {
		m.connections = connections
	}
}

func WithWeight(weight int) MetricsOption {
	if weight == 0 {
		weight = 1
	}
	return func(m *Metrics) {
		m.weight = weight
	}
}
