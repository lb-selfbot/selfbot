//go:build lbmobile

package main

import (
	"context"
	"net"
	"runtime/debug"

	"github.com/LightningDev1/LB-Selfbot/mobile"
)

func main() {
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(context, "udp", "1.1.1.1:53")
		},
	}

	debug.SetMemoryLimit(300 * 1 << 20)

	mobile.Main()
}
