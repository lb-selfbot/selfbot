package network

import (
	"net"
	"net/url"
	"time"

	"github.com/lb-selfbot/selfbot/internal/http"
)

// GetPingLatency pings a host and returns the round-trip time
func GetPingLatency(host string) (time.Duration, error) {
	if url, err := url.Parse(host); err == nil {
		host = url.Hostname()
	}

	host = net.JoinHostPort(host, "80")

	start := time.Now()
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return time.Since(start), nil
}

// GetPublicIP retrieves the public IP address using an external service
func GetPublicIP() string {
	resp, err := http.DefaultClient().Get("https://api.ipify.org")
	if err != nil {
		return "Error"
	}

	return string(resp.MustGetBody())
}

// GetLocalIP retrieves the local IP address on the network
func GetLocalIP() string {
	// Get preferred outbound IP address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetMacAddress retrieves the MAC address of the first non-loopback network interface
func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Unknown"
	}

	for _, iface := range interfaces {
		// Skip loopback and interfaces without MAC
		if iface.Flags&net.FlagLoopback != 0 || iface.HardwareAddr == nil {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" {
			return mac
		}
	}

	return "Unknown"
}
