package dnsutils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"unicode/utf8"

	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

var (
	FetchIPv4AddrUrl string
	FetchIPv6AddrUrl string
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
	httpClient := &http.Client{}

	var ipv6URL = "http://api6.ipify.org"
	if FetchIPv6AddrUrl != "" {
		ipv6URL = FetchIPv6AddrUrl
	}
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
	httpClient := &http.Client{}

	var ipv4URL = "http://api.ipify.org"
	if FetchIPv4AddrUrl != "" {
		ipv4URL = FetchIPv4AddrUrl
	}
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

func convrtToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return string(strBytes)
}

func DoCommand(cmd string, args string) string {
	out, err := exec.Command(cmd, args).Output()
	if err != nil {
		return ""
	}

	if !utf8.Valid(out) {
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest(out)
		var encoding string
		if err == nil {
			if result.Language == "zh" {
				encoding = "gbk"
			} else {
				encoding = result.Charset
			}
		}
		utf8 := convrtToUTF8(string(out), encoding)
		if err != nil {
			return string(out)
		}
		return string(utf8)
	}

	return string(out)
}
