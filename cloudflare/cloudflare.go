package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
)

// ipv6RegisterToCF ...
func ipv6RegisterToCF(domain, token, zone_id, rid, ip string) error {
	url := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records/" + rid

	bodyMap := map[string]interface{}{
		"type":    "AAAA",
		"name":    domain,
		"content": ip,
		"ttl":     60,
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	// 使用 PUT 方法更新记录.
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(bodyBytes)))
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("cloudflare API returned status %d: %s", res.StatusCode, string(body))
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("cloudflare API error: %s", string(body))
	}

	return nil
}

// ipv4RegisterToCF ...
func ipv4RegisterToCF(domain, token, zone_id, rid, ip string) error {
	url := "https://api.cloudflare.com/client/v4/zones/" + zone_id + "/dns_records/" + rid

	bodyMap := map[string]interface{}{
		"type":    "A",
		"name":    domain,
		"content": ip,
		"ttl":     60,
	}

	bodyBytes, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	// 使用 PUT 方法更新记录.
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(bodyBytes)))
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("cloudflare API returned status %d: %s", res.StatusCode, string(body))
	}

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf("cloudflare API error: %s", string(body))
	}

	return nil
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

// FetchZoneID ...
func FetchZoneID(domain, token string) (string, error) {
	url := "https://api.cloudflare.com/client/v4/zones?name=" + domain

	req, err := http.NewRequest("GET", url, nil)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("cloudflare API returned status %d: %s", res.StatusCode, string(body))
	}

	var result struct {
		Success bool                   `json:"success"`
		Result  []map[string]interface{} `json:"result"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if !result.Success {
		return "", errors.New(string(body))
	}

	if len(result.Result) == 0 {
		return "", fmt.Errorf("zone not found for domain: %s", domain)
	}

	// 获取zone_id.
	zoneID, ok := result.Result[0]["id"].(string)
	if !ok {
		return "", fmt.Errorf("unable to parse zone_id from response: %s", string(body))
	}

	return zoneID, nil
}

// FetchRecordID ...
func FetchRecordID(zone_id, token, domain, dnsType string) (string, error) {
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("cloudflare API returned status %d: %s", res.StatusCode, string(body))
	}

	var result cfResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if !result.Success {
		return "", fmt.Errorf("cloudflare API error: %s", string(body))
	}

	for _, element := range result.Records {
		if element.Name == domain && element.Type == dnsType {
			return element.ID, nil
		}
	}

	return "", fmt.Errorf("DNS record not found: %s (%s) in zone %s", domain, dnsType, zone_id)
}
