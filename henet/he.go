package henet

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

// 下面为 he.net 的 API 说明

// 自动检测我的 IPv4/IPv6 地址：
// % curl -4 "http://dyn.example.com:password@dyn.dns.he.net/nic/update?hostname=dyn.example.com"
// % curl -6 "http://dyn.example.com:password@dyn.dns.he.net/nic/update?hostname=dyn.example.com"

// 手动设置 IPv4/IPv6 地址：
// % curl "http://dyn.example.com:password@dyn.dns.he.net/nic/update?hostname=dyn.example.com&myip=192.168.0.1"
// % curl "http://dyn.example.com:password@dyn.dns.he.net/nic/update?hostname=dyn.example.com&myip=2001:db8:beef:cafe::1"

// 注意：用户名也是主机名。密码使用“password=”发送。这会跳过 HTTP 基本身份验证。
// 使用 GET 进行身份验证和更新：
// % curl "https://dyn.dns.he.net/nic/update?hostname=dyn.example.com&password=password&myip=192.168.0.1"
// % curl "https://dyn.dns.he.net/nic/update?hostname=dyn.example.com&password=password&myip=2001:db8:beef:cafe::1"

// 使用 POST 进行身份验证和更新：
// % curl "https://dyn.dns.he.net/nic/update" -d "hostname=dyn.example.com" -d "password=password" -d "myip=192.168.0.1"
// % curl "https://dyn.dns.he.net/nic/update" -d "hostname=dyn.example.com" -d "password=password" -d "myip=2001:db8:beef:cafe::1"

func registerToHenet(domain, subdomain, passwd, ip string) error {
	domain = subdomain + "." + domain
	Url := "https://dyn.dns.he.net/nic/update?" +
		"hostname=" + domain +
		"&password=" + passwd +
		"&myip=" + ip

	request, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return err
	}

	henetClient := &http.Client{}
	res, err := henetClient.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return err
}

// DoHenetv6 ...
func DoHenetv6(domain, subdomain, passwd, extIP string) {
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

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向godaddy等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv6address")
	if err == nil {
		buf := make([]byte, 1024)
		f.Read(buf)
		f.Close()

		// 获取ip字符串.
		storeIP = strings.TrimRight(string(buf), string(rune(0)))
	}

	if storeIP == ipv6 {
		info := "ipv6 " + storeIP + " same as " + ipv6
		fmt.Println(info)
		return
	}

	err = registerToHenet(domain, subdomain, passwd, ipv6)
	if err != nil {
		fmt.Println("register to he.net error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoHenetv4 ...
func DoHenetv4(domain, subdomain, ssoKey, extIP string) {
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

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向godaddy等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv4address")
	if err == nil {
		buf := make([]byte, 1024)
		f.Read(buf)
		f.Close()

		// 获取ip字符串.
		storeIP = strings.TrimRight(string(buf), string(rune(0)))
	}

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = registerToHenet(domain, subdomain, ssoKey, ipv4)
	if err != nil {
		fmt.Println("register to he.net error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
}
