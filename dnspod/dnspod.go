package dnspod

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

// ipv6RegisterToDNSPOD ...
func ipv6RegisterToDNSPOD(domain, subdomain, token, rid, ip string) error {
	Url := "https://dnsapi.cn/Record.Modify"

	res, err := http.PostForm(Url, url.Values{
		"login_token": {token},
		"format":      {"json"},
		"domain":      {domain},
		"record_id":   {rid},
		"sub_domain":  {subdomain},
		"record_type": {"AAAA"},
		"record_line": {"默认"},
		"value":       {ip},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// ipv4RegisterToDNSPOD ...
func ipv4RegisterToDNSPOD(domain, subdomain, token, rid, ip string) error {
	Url := "https://dnsapi.cn/Record.Modify"

	res, err := http.PostForm(Url, url.Values{
		"login_token": {token},
		"format":      {"json"},
		"domain":      {domain},
		"record_id":   {rid},
		"sub_domain":  {subdomain},
		"record_type": {"A"},
		"record_line": {"默认"},
		"value":       {ip},
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// DoDNSPODv6 ...
func DoDNSPODv6(domain, subdomain, token, rid, extIP string) {
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
	// 否则向dnspod等域名服务注册修改ip, 并保存ip
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

	err = ipv6RegisterToDNSPOD(domain, subdomain, token, rid, ipv6)
	if err != nil {
		fmt.Println("register to dnspod error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoDNSPODv4 ...
func DoDNSPODv4(domain, subdomain, token, rid, extIP string) {
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

	// 获取ip字符串.
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

		storeIP = strings.TrimRight(string(buf), string(rune(0)))
	}

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = ipv4RegisterToDNSPOD(domain, subdomain, token, rid, ipv4)
	if err != nil {
		fmt.Println("register to dnspod error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
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

// FetchRecordID ...
func FetchRecordID(token, domain, subdomain, domainType string) (string, error) {
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
