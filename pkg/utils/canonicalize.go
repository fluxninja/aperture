package utils

import "strings"

// CanonicalizeOtelHeaderKey converts header naming convention to Otel's one.
func CanonicalizeOtelHeaderKey(key string) string {
	return strings.ReplaceAll(strings.ToLower(key), "-", "_")
}

// CanonicalizeOtelHTTPFlavor converts protocol to Otel kind of HTTP protocol.
func CanonicalizeOtelHTTPFlavor(protocolName string) string {
	switch protocolName {
	case "HTTP/1.0":
		return "1.0"
	case "HTTP/1.1":
		return "1.1"
	case "HTTP/2":
		return "2.0"
	case "HTTP/3":
		return "3.0"
	default:
		return protocolName
	}
}
