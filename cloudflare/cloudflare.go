package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

// ipv6RegisterToCF ...
func ipv6RegisterToCF(domain, token, zone_id, rid, ip string) error {
	Url := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records/" + rid

	// 使用 PUT 方法更新记录.
	req, err := http.NewRequest("PUT", Url,
		strings.NewReader(`{"type":"AAAA","name":"`+domain+`","content":"`+ip+`","ttl":60}`))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// ipv4RegisterToCF ...
func ipv4RegisterToCF(domain, token, zone_id, rid, ip string) error {
	Url := "https://api.cloudflare.com/client/v4/zones/" + domain + "/dns_records/" + rid

	// 使用 PUT 方法更新记录.
	req, err := http.NewRequest("PUT", Url,
		strings.NewReader(`{"type":"A","name":"`+domain+`","content":"`+ip+`","ttl":60}`))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	return err
}

// DoCFv6 ...
func DoCFv6(domain, token, zone_id, rid, extIP string) {
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
	// 否则向cloudflare等域名服务注册修改ip, 并保存ip
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

	err = ipv6RegisterToCF(domain, token, zone_id, rid, ipv6)
	if err != nil {
		fmt.Println("register to cloudflare error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoCFv4 ...
func DoCFv4(domain, token, zone_id, rid, extIP string) {
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

	err = ipv4RegisterToCF(domain, token, zone_id, rid, ipv4)
	if err != nil {
		fmt.Println("register to cloudflare error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
}

type cfRecordResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}

type cfResult struct {
	Success bool             `json:"success"`
	Records []cfRecordResult `json:"result"`
}

// GetZoneID ...
func FetchZoneID(domain, token string) (string, error) {
	Url := "https://api.cloudflare.com/client/v4/zones?name=" + domain

	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result["success"] == false {
		return "", errors.New(string(body))
	}

	// 获取zone_id.
	zone_id := result["result"].([]interface{})[0].(map[string]interface{})["id"].(string)

	return zone_id, nil
}

// FetchRecordID ...
func FetchRecordID(zone_id, token, domain string) (string, error) {
	URL := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records"

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result cfResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	for _, element := range result.Records {
		if element.Name == domain {
			return element.ID, nil
		}
	}

	return "", errors.New(string(body))
}
