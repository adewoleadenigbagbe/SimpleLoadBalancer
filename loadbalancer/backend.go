package loadbalancer

import (
	"database/sql"
	"net/http/httputil"
	"net/url"
	"sync"
)

var _ IBackend = (*Backend)(nil)

type IBackend interface {
	GetID() string
	SetAlive(bool)
	IsAlive() bool
	GetURL() url.URL
}

type Metrics struct {
	connections sql.NullInt32
	weight      sql.NullFloat64
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
	backend.alive = alive
}

func (backend *Backend) GetURL() url.URL {
	return backend.url
}

func (backend *Backend) IsAlive() bool {
	return backend.alive
}
