package externalip

import (
	"net"
	"encoding/base64"
)

// IsIPv6 ...
func IsIPv6(str string) bool {
        ip := net.ParseIP(str)
        return ip.To4() == nil
}

// IsIPv4 ...
func IsIPv4(str string) bool {
        ip := net.ParseIP(str)
        return ip.To4() != nil
}

// BasicAuth ...
func BasicAuth(username, password string) string {
        auth := username + ":" + password
        return base64.StdEncoding.EncodeToString([]byte(auth))
}



