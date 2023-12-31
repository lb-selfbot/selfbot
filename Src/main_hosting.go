//go:build hosting

package main

import (
	"context"
	"net"

	"github.com/LightningDev1/LB-Selfbot/hosting"
)

func init() {
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(context, "udp", "1.1.1.1:53")
		},
	}
}

func main() {
	hosting.Main()
}
