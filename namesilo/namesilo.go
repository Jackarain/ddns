package namesilo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

type ResourceRecord struct {
	RecordId string `xml:"record_id"`
	Type     string `xml:"type"`
	Host     string `xml:"host"`
	Value    string `xml:"value"`
	Ttl      int    `xml:"ttl"`
}

type NamesiloReply struct {
	ResourceRecords []ResourceRecord `xml:"resource_record"`
}

type Namesilo struct {
	Reply NamesiloReply `xml:"reply"`
}

// FetchRecordID ...
func FetchRecordID(token, domain, subdomain string) (string, error) {
	Url := fmt.Sprintf("https://www.namesilo.com/api/dnsListRecords?version=1&type=xml&key=%s&domain=%s", token, domain)

	request, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return "", err
	}

	namesiloClient := &http.Client{}
	res, err := namesiloClient.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var namesilo Namesilo
	err = xml.Unmarshal(body, &namesilo)
	if err != nil {
		fmt.Println("XML unmarshal error:", err)
		return "", err
	}

	hostname := subdomain + "." + domain
	for _, element := range namesilo.Reply.ResourceRecords {
		if element.Host == hostname {
			return element.RecordId, nil
		}
	}

	return "", errors.New("not found record id")
}

func registerIP(domain, subdomain, token, rid, extIP string) error {
	namesiloURL := fmt.Sprintf("https://www.namesilo.com/api/dnsUpdateRecord?version=1&type=xml&key=%s&domain=%s&rrid=%s&rrhost=%s&rrvalue=%s",
		token, domain, rid, subdomain, extIP)

	request, err := http.NewRequest("GET", namesiloURL, nil)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}

	namesiloClient := &http.Client{}
	res, err := namesiloClient.Do(request)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return nil
}

// DoNamesiloV6 ...
func DoNamesiloV6(domain, subdomain, token, rid, extIP string) {
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

	err = registerIP(domain, subdomain, token, rid, ipv6)
	if err != nil {
		fmt.Println("register to namesilo error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoNamesiloV4 ...
func DoNamesiloV4(domain, subdomain, token, rid, extIP string) {
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

	err = registerIP(domain, subdomain, token, rid, ipv4)
	if err != nil {
		fmt.Println("register to namesilo error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
}
