package backend

import (
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBackend(t *testing.T) {
	url, _ := url.Parse("http://localhost:3333")
	b := NewBackend(url, httputil.NewSingleHostReverseProxy(url)).(*Backend)
	assert.Equal(t, "http://localhost:3333", b.url.String())
}

func TestSetAlive(t *testing.T) {
	url, _ := url.Parse("http://localhost:3333")
	b := NewBackend(url, httputil.NewSingleHostReverseProxy(url)).(*Backend)
	b.SetAlive(true)
	assert.Equal(t, true, b.alive)
}
