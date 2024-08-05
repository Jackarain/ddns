package alidns

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Jackarain/ddns/dnsutils"
)

var (
	User   string
	Passwd string
)

// generateSignature generates a signature for the API request
func generateSignature(params map[string]string, secret string) string {
	// Step 1: Sort the parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Step 2: Concatenate the parameters
	var sortedParams string
	for _, k := range keys {
		sortedParams += "&" + specialURLEncode(k) + "=" + specialURLEncode(params[k])
	}
	stringToSign := "GET&%2F&" + specialURLEncode(sortedParams[1:])

	// Step 3: Calculate HMAC SHA1
	h := hmac.New(sha1.New, []byte(secret+"&"))
	h.Write([]byte(stringToSign))

	// Step 4: Base64 encode the result
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// specialURLEncode encodes a string for URL, following specific rules
func specialURLEncode(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}

// sendRequest sends a GET request to the API and returns the response body
func sendRequest(params map[string]string) ([]byte, error) {
	// Add common parameters
	params["Format"] = "JSON"
	params["Version"] = "2015-01-09"
	params["AccessKeyId"] = User
	params["SignatureMethod"] = "HMAC-SHA1"
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	params["SignatureVersion"] = "1.0"
	params["SignatureNonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	// Generate the signature
	params["Signature"] = generateSignature(params, Passwd)

	// Construct the URL
	var query string
	for k, v := range params {
		query += "&" + k + "=" + v
	}
	url := "https://alidns.aliyuncs.com/?" + query[1:]

	// Send the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// registerToAlidns updates a DNS record in AliDNS
func registerToAlidns(domain, subdomain, rid, ip, tp string) error {
	params := map[string]string{
		"Action":   "UpdateDomainRecord",
		"RecordId": rid,
		"RR":       subdomain,
		"Type":     tp,
		"Value":    ip,
		"TTL":      "600",
	}

	_, err := sendRequest(params)
	return err
}

// DoAlidnsV6 handles IPv6 DNS updates
func DoAlidnsV6(domain, subdomain, passwd, extIP string) {
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

	err = registerToAlidns(domain, subdomain, passwd, ipv6, "AAAA")
	if err != nil {
		fmt.Println("register to alidns error: ", err)
		return
	}

	dnsutils.FileWriteString("ipv6address", ipv6)
}

// DoAlidnsV4 handles IPv4 DNS updates
func DoAlidnsV4(domain, subdomain, rid, extIP string) {

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

	err = registerToAlidns(domain, subdomain, rid, ipv4, "A")
	if err != nil {
		fmt.Println("register to alidns error: ", err)
		return
	}

	dnsutils.FileWriteString("ipv4address", ipv4)
}

// FetchRecordID fetches the record ID for a domain
func FetchRecordID(domain string) (string, error) {
	params := map[string]string{
		"Action":     "DescribeDomainRecords",
		"DomainName": domain,
	}

	respBody, err := sendRequest(params)
	if err != nil {
		return "", err
	}

	return parseRecordID(respBody)
}

// parseRecordID extracts the Record ID from the JSON response
func parseRecordID(respBody []byte) (string, error) {
	// JSON 内容参考：
	// {
	//   "TotalCount": 2,
	//   "PageSize": 20,
	//   "RequestId": "536E9CAD-DB30-4647-AC87-AA5CC38C5382",
	//   "DomainRecords": {
	//     "Record": [
	//       {
	//         "Status": "Enable",
	//         "Type": "MX",
	//         "Remark": "备注",
	//         "TTL": 600,
	//         "RecordId": "9999985",
	//         "Priority": 5,
	//         "RR": "www",
	//         "DomainName": "example.com",
	//         "Weight": 2,
	//         "Value": "mail1.hichina.com",
	//         "Line": "default",
	//         "Locked": false,
	//         "CreateTimestamp": 1666501957000,
	//         "UpdateTimestamp": 1676872961000
	//       }
	//     ]
	//   },
	//   "PageNumber": 1
	// }

	// 定义一个结构来解析 JSON
	type Record struct {
		RecordId string `json:"RecordId"`
	}

	type DomainRecords struct {
		Record []Record `json:"Record"`
	}

	type Response struct {
		DomainRecords DomainRecords `json:"DomainRecords"`
	}

	var resp Response

	// 解析 JSON
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", err
	}

	// 检查是否有 Record
	if len(resp.DomainRecords.Record) == 0 {
		return "", errors.New("no records found")
	}

	// 返回第一个 Record 的 Record ID
	return resp.DomainRecords.Record[0].RecordId, nil
}
