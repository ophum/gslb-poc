package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second,
			}
			return d.DialContext(ctx, "udp", "127.0.0.1:8053")
		},
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := net.Dialer{
			Timeout:  time.Second,
			Resolver: &resolver,
		}
		return d.DialContext(ctx, network, addr)
	}

	resp, err := http.Get("http://www.example.com:8000")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("host:", resp.Request.Host)
	fmt.Println("body:", string(body))
}
