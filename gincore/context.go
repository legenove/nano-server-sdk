/*
	for prepare context
*/
package gincore

import (
	"net"
	"net/http"
	"strings"
)

func RequestIP(req *http.Request) string {
	if req == nil {
		return ""
	}
	var ip string
	ip = req.Header.Get("X-Forwarded-For")
	if ip != "" {
		i := strings.Index(ip, ",")
		if i != -1 {
			ip = ip[:i]
		}
		return ip
	}
	ip = req.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}
