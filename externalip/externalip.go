package externalip

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/axgle/mahonia"
)

// FileWriteString ...
func FileWriteString(name, ip string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.WriteString(ip)
}

// FileReadString ...
func FileReadString(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

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

// ExternalIPv6 ...
func ExternalIPv6() (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: false,
	}
	httpClient := &http.Client{
		Transport: tr,
	}

	const ipv6URL = "https://api6.ipify.org"
	req, err := http.NewRequest("GET", ipv6URL, nil)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ipv6 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	return ipv6, err
}

// ExternalIPv4 ...
func ExternalIPv4() (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: false,
	}
	httpClient := &http.Client{
		Transport: tr,
	}

	const ipv4URL = "https://api.ipify.org"
	req, err := http.NewRequest("GET", ipv4URL, nil)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ipv4 := mahonia.NewDecoder("gbk").ConvertString(string(body))
	return ipv4, err
}
