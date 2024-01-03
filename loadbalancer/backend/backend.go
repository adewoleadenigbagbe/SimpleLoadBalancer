package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

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
	GetResponseTime() int64
	Serve(w http.ResponseWriter, r *http.Request)
}

type Metrics struct {
	connections  int
	weight       int
	responseTime time.Duration
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
	alive := backend.alive
	defer backend.mux.RUnlock()
	return alive
}

func (backend *Backend) GetActiveConnections() int {
	backend.mux.RLock()
	connections := backend.metrics.connections
	defer backend.mux.RUnlock()
	return connections
}

func (backend *Backend) GetWeight() int {
	backend.mux.RLock()
	weight := backend.metrics.weight
	defer backend.mux.RUnlock()
	return weight
}

func (backend *Backend) GetResponseTime() int64 {
	return backend.metrics.responseTime.Microseconds()
}

func (backend *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	beforeServe := time.Now()
	backend.reverseProxy.ServeHTTP(w, r)
	afterServe := time.Now()

	backend.mux.Lock()
	backend.metrics.connections++
	duration := afterServe.Sub(beforeServe)
	backend.metrics.responseTime += duration
	backend.mux.Unlock()
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

func WithResponseTime(d time.Duration) MetricsOption {
	return func(m *Metrics) {
		m.responseTime = d
	}
}
