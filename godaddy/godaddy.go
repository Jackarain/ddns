package godaddy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"rpi4.p2sp.net/Jackarain/ddns/externalip"
)

// IPv6RegisterToGodaddy ...
func IPv6RegisterToGodaddy(domain, subdomain, ssoKey, ip string) error {
	URL := "https://api.godaddy.com/v1/domains/" + domain + "/records/AAAA/" + subdomain
	payload := strings.NewReader("[{\"data\": \"" + ip + "\"}]")
	req, _ := http.NewRequest("PUT", URL, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "sso-key "+ssoKey)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err == nil {
		return nil
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return err
}

// IPv4RegisterToGodaddy ...
func IPv4RegisterToGodaddy(domain, subdomain, ssoKey, ip string) error {
	URL := "https://api.godaddy.com/v1/domains/" + domain + "/records/A/" + subdomain
	payload := strings.NewReader("[{\"data\": \"" + ip + "\"}]")
	req, _ := http.NewRequest("PUT", URL, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "sso-key "+ssoKey)

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err == nil {
		return nil
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return err
}

// DoGodaddyv6 ...
func DoGodaddyv6(domain, subdomain, ssoKey string) {
	// 获取公网IPv6地址.
	ipv6, err := externalip.ExternalIPv6()
	if err != nil {
		fmt.Println("ipv6: ", err)
		return
	}
	if !externalip.IsIPv6(ipv6) {
		fmt.Println("external ipv6 error:", ipv6)
		return
	}

	fmt.Println("external ipv6: ", ipv6)

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向godaddy等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv6address")
	if err != nil {
		externalip.FileWriteString("ipv6address", ipv6)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		externalip.FileWriteString("ipv6address", ipv6)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv6 {
		info := "ipv6 " + storeIP + " same as " + ipv6
		fmt.Println(info)
		return
	}

	err = IPv6RegisterToGodaddy(domain, subdomain, ssoKey, ipv6)
	if err != nil {
		fmt.Println("register to godaddy error: ", err)
		return
	}

	// 重写ip缓存文件.
	externalip.FileWriteString("ipv6address", ipv6)
}

// DoGodaddyv4 ...
func DoGodaddyv4(domain, subdomain, ssoKey string) {
	// 获取公网IPv4地址.
	ipv4, err := externalip.ExternalIPv4()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("external ipv4: ", ipv4)

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向godaddy等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv4address")
	if err != nil {
		externalip.FileWriteString("ipv4address", ipv4)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		externalip.FileWriteString("ipv4address", ipv4)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = IPv4RegisterToGodaddy(domain, subdomain, ssoKey, ipv4)
	if err != nil {
		fmt.Println("register to godaddy error: ", err)
		return
	}

	// 重写ip缓存文件.
	externalip.FileWriteString("ipv4address", ipv4)
}
