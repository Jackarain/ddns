package oray

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

var (
	User   string
	Passwd string
)

// http://ddns.oray.com/ph/update?hostname=660kw424.goho.co&myip=1.1.1.1

// registerToOray ...
func registerToOray(domain, ip string) error {
	Url := "http://ddns.oray.com/ph/update?" +
		"hostname=" + domain +
		"&myip=" + ip

	request, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Basic "+dnsutils.BasicAuth(User, Passwd))

	orayClient := &http.Client{}
	res, err := orayClient.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// DoOrayv4 ...
func DoOrayv4(domain, extIP string) {
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
	// 否则向域名服务注册修改ip, 并保存ip
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

	err = registerToOray(domain, ipv4)
	if err != nil {
		fmt.Println("register to f3322 error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
}
