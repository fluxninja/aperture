package middlewares

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/go-logr/logr"
)

// splits address into host and port
func splitAddress(logger logr.Logger, address string) (string, uint32) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		logger.V(2).Info("Error splitting address", "address", address, "error", err)
		return host, 0
	}

	portUint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		logger.V(2).Info("Error parsing port", "address", address, "error", err)
		return host, 0
	}

	return host, uint32(portUint)
}

// reads body from request, replacing it with a clone to allow further reads
func readClonedBody(r *http.Request) ([]byte, error) {
	body := r.Body
	defer body.Close()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// try to figure out the local ip address
func getLocalIP(logger logr.Logger) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.V(2).Info("Failed to get local IP address", "error", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	logger.V(2).Info("Failed to get local IP address")
	return ""
}
