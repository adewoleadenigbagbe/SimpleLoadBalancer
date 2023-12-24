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
	SetAlive(bool)
	IsAlive() bool
	GetURL() url.URL
	Serve(w http.ResponseWriter, r *http.Request)
}

type Metrics struct {
	//connections sql.NullInt32
	//weight      sql.NullFloat64

	connections int
	weight      float64
}

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

func NewBackend(endpoint *url.URL, proxy *httputil.ReverseProxy) IBackend {
	backend := Backend{
		id:           uuid.NewString(),
		url:          *endpoint,
		reverseProxy: proxy,
		metrics:      &Metrics{},
	}

	return &backend
}
