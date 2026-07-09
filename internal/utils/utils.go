package utils

import (
	"net"
	"os"
	"strings"
)

func IsGoRun() bool {
	execPath := os.Args[0]
	return strings.Contains(execPath, "go-build")
}

func GetIps() []string {
	ips := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() != nil && (ip.IsPrivate() || ip.IsLoopback()) {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips
}
