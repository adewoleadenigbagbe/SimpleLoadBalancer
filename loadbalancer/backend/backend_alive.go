package backend

import (
	"context"
	"fmt"
	"net"
	"net/url"
)

func IsBackendAlive(ctx context.Context, aliveChannel chan bool, u url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		fmt.Println("Site Unreachable")
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}
