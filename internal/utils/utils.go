package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"runtime"
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

func Jn(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprint(&sb, arg)
	}
	return sb.String()
}

func OpenBrowser(url string) error {

	if IsGoRun() {
		return nil
	}

	var cmd string
	args := []string{}

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = append(args, "/c", "start")
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)

	c := exec.Command(cmd, args...)

	err := c.Start()
	if err != nil {
		slog.Error("Failed to open browser", "err", err)
		return err
	}

	return c.Wait()
}

// Deprecated: Very slow, use only for tests.
func DeepCopy[T any](src T) (T, error) {
	var buf bytes.Buffer
	var dst T

	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return dst, err
	}
	if err := gob.NewDecoder(&buf).Decode(&dst); err != nil {
		return dst, err
	}
	return dst, nil
}
