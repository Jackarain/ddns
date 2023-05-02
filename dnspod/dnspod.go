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

// RegisterToF3322 ...
func RegisterToF3322(ip string) error {
	f3322Url := "http://members.3322.net/dyndns/update?system=dyndns&hostname=sgrc.f3322.net&myip=" + ip
	request, err := http.NewRequest("GET", f3322Url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Basic "+dnsutils.BasicAuth("wgm001", "ggfggc"))
	f3322Client := &http.Client{}
	response, err := f3322Client.Do(request)
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

// IPv6RegisterToDNSPOD ...
func IPv6RegisterToDNSPOD(domain, subdomain, rid, ip string) error {
	dnspodURL := "https://dnsapi.cn/Record.Modify"
	response, err := http.PostForm(dnspodURL, url.Values{
		"login_token": {"11898,8d7347cd5969f7aa89752c068a6b949a"},
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
	response.Body.Close()
	return nil
}

// IPv4RegisterToDNSPOD ...
func IPv4RegisterToDNSPOD(domain, subdomain, rid, ip string) error {
	dnspodURL := "https://dnsapi.cn/Record.Modify"
	response, err := http.PostForm(dnspodURL, url.Values{
		"login_token": {"11898,8d7347cd5969f7aa89752c068a6b949a"},
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
	response.Body.Close()
	return nil
}

// DoDNSPODv6 ...
func DoDNSPODv6(domain, subdomain, rid, extIP string) {
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

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向dnspod等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv6address")
	if err != nil {
		dnsutils.FileWriteString("ipv6address", ipv6)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		dnsutils.FileWriteString("ipv6address", ipv6)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv6 {
		info := "ipv6 " + storeIP + " same as " + ipv6
		fmt.Println(info)
		return
	}

	err = IPv6RegisterToDNSPOD(domain, subdomain, rid, ipv6)
	if err != nil {
		fmt.Println("register to dnspod error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoDNSPODv4 ...
func DoDNSPODv4(domain, subdomain, rid, extIP string) {
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

	fmt.Println("external ipv4: ", ipv4)

	// 如果能打开ipaddress, 则读取ipaddress中的ip
	// 与获取的公网ip对比, 如果没有改变, 则退出,
	// 否则向dnspod等域名服务注册修改ip, 并保存ip
	// 到文件 ipaddress 中.
	f, err := os.Open("ipv4address")
	if err != nil {
		dnsutils.FileWriteString("ipv4address", ipv4)
	}

	buf := make([]byte, 1024)
	n, _ := f.Read(buf)
	if n == 0 {
		dnsutils.FileWriteString("ipv4address", ipv4)
	}
	f.Close()

	// 获取ip字符串.
	storeIP := strings.TrimRight(string(buf), string(0))

	if storeIP == ipv4 {
		info := "ipv4 " + storeIP + " same as " + ipv4
		fmt.Println(info)
		return
	}

	err = IPv4RegisterToDNSPOD(domain, subdomain, rid, ipv4)
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
