package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/axgle/mahonia"
	"rpi4.p2sp.net/Jackarain/externalip"
	"rpi4.p2sp.net/Jackarain/ddns/ddnsreg"
)

func externalIPv6() (string, error) {
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

func externalIPv4() (string, error) {
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

func fileWriteString(name, ip string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.WriteString(ip)
}

func fileReadString(name string) (string, error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func doDNSPODv6(domain, subdomain, rid string) {
	// 获取公网IPv6地址.
	ipv6, err := externalIPv6()
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
	// 否则向dnspod等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv6address")
	if err != nil {
		fileWriteString("ipv6address", ipv6)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		fileWriteString("ipv6address", ipv6)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv6 {
		info := "ipv6 " + storeIP + " same as " + ipv6
		fmt.Println(info)
		return
	}

	err = ddnsreg.IPv6RegisterToDNSPOD(domain, subdomain, rid, ipv6)
	if err != nil {
		fmt.Println("register to dnspod error: ", err)
		return
	}

	// 重写ip缓存文件.
	fileWriteString("ipv6address", ipv6)
}

func doDNSPODv4(domain, subdomain, rid string) {
	// 获取公网IPv4地址.
	ipv4, err := externalIPv4()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("external ipv4: ", ipv4)

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向dnspod等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv4address")
	if err != nil {
		fileWriteString("ipv4address", ipv4)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		fileWriteString("ipv4address", ipv4)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = ddnsreg.IPv4RegisterToDNSPOD(domain, subdomain, rid, ipv4)
	if err != nil {
		fmt.Println("register to dnspod error: ", err)
		return
	}

	// 重写ip缓存文件.
	fileWriteString("ipv4address", ipv4)
}

type dnspodStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type dnspodRecordInfo struct {
	ID    string `json:"id"`
	Line  string `json:"line"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type dnspodResult struct {
	Status  dnspodStatus       `json:"status"`
	Records []dnspodRecordInfo `json:"records"`
}

func fetchRecordID(token, domain, subdomain, domainType string) (string, error) {
	dnspodURL := "https://dnsapi.cn/Record.List"

	response, err := http.PostForm(dnspodURL, url.Values{
		"login_token": {token},
		"format":      {"json"},
		"domain":      {domain},
		"subdomain":   {subdomain},
		"length":      {"3000"},
		"record_type": {domainType},
	})
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var result dnspodResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	for _, element := range result.Records {
		if element.Name == subdomain {
			return element.ID, nil
		}
	}

	return "", errors.New(result.Status.Message)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: ./ddns <token> <domain> <subdomain> [type]")
		return
	}

	if len(args) == 4 {
		ridFileName := args[2] + args[3]
		rid, err := fileReadString(ridFileName)
		if err != nil || rid == "" {
			rid, err = fetchRecordID(args[0], args[1], args[2], args[3])
			if err != nil {
				fmt.Println(err)
				return
			}

			fileWriteString(ridFileName, rid)
		}

		if args[3] == "A" {
			doDNSPODv4(args[1], args[2], rid)
		} else if args[3] == "AAAA" {
			doDNSPODv6(args[1], args[2], rid)
		} else {
			fmt.Println(err)
		}

		return
	}

	ridFileName := args[2] + "A"
	rid, err := fileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = fetchRecordID(args[0], args[1], args[2], "A")
		if err != nil {
			fmt.Println(err)
			return
		}
		fileWriteString(ridFileName, rid)
	}

	fmt.Println(args[2], "a record id:", rid)
	doDNSPODv4(args[1], args[2], rid)

	ridFileName = args[2] + "AAAA"
	rid, err = fileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = fetchRecordID(args[0], args[1], args[2], "AAAA")
		if err != nil {
			fmt.Println(err)
			return
		}
		fileWriteString(ridFileName, rid)
	}

	fmt.Println(args[2], "aaaa record id:", rid)
	doDNSPODv6(args[1], args[2], rid)
}
