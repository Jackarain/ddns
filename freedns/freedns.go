package freedns

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

// FreeDNS API:
// 更新动态DNS记录(IPv4):
//   https://sync.afraid.org/u/TOKEN/?myip=IP_ADDRESS
// 更新动态DNS记录(IPv6):
//   http://v6.sync.afraid.org/u/TOKEN/?myip=IP_ADDRESS
// 响应:
//   "Updated X.X.X.X"  - 成功更新
//   "No change"        - IP未发生变化
//   "ERROR: ..."       - 发生错误

func registerToFreeDNSv4(token, ip string) error {
	url := "https://sync.afraid.org/u/" + token + "/?myip=" + ip

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

func registerToFreeDNSv6(token, ip string) error {
	url := "http://v6.sync.afraid.org/u/" + token + "/?myip=" + ip

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// DoFreeDNSv4 ...
func DoFreeDNSv4(token, extIP string) {
	var ipv4 string
	if extIP == "" {
		ip, err := dnsutils.ExternalIPv4()
		if err != nil {
			fmt.Println(err)
			return
		}
		ipv4 = ip
	} else {
		ipv4 = extIP
	}

	if len(ipv4) == 0 {
		return
	}

	fmt.Println("external ipv4: ", ipv4)

	var storeIP string

	f, err := os.Open("ipv4address")
	if err == nil {
		buf := make([]byte, 1024)
		f.Read(buf)
		f.Close()

		storeIP = strings.TrimRight(string(buf), string(rune(0)))
	}

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = registerToFreeDNSv4(token, ipv4)
	if err != nil {
		fmt.Println("register to freedns error: ", err)
		return
	}

	dnsutils.FileWriteString("ipv4address", ipv4)
}

// DoFreeDNSv6 ...
func DoFreeDNSv6(token, extIP string) {
	var ipv6 string
	if extIP == "" {
		ip, err := dnsutils.ExternalIPv6()
		if err != nil {
			fmt.Println("ipv6: ", err)
			return
		}
		ipv6 = ip
	} else {
		ipv6 = extIP
	}

	if !dnsutils.IsIPv6(ipv6) {
		fmt.Println("external ipv6 error:", ipv6)
		return
	}
	fmt.Println("external ipv6: ", ipv6)

	var storeIP string

	f, err := os.Open("ipv6address")
	if err == nil {
		buf := make([]byte, 1024)
		f.Read(buf)
		f.Close()

		storeIP = strings.TrimRight(string(buf), string(rune(0)))
	}

	if storeIP == ipv6 {
		info := "ipv6 " + storeIP + " same as " + ipv6
		fmt.Println(info)
		return
	}

	err = registerToFreeDNSv6(token, ipv6)
	if err != nil {
		fmt.Println("register to freedns error: ", err)
		return
	}

	dnsutils.FileWriteString("ipv6address", ipv6)
}
