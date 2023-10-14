package alidns

import (
	"fmt"
	"os"
	"strings"

	"github.com/Jackarain/ddns/dnsutils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

var (
	User   string
	Passwd string
)

// registerToAlidns ...
func registerToAlidns(domain, subdomain, rid, ip string) error {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", User, Passwd)
	if err != nil {
		return err
	}

	req := alidns.CreateUpdateDomainRecordRequest()
	req.Scheme = "https"
	req.RecordId = rid
	req.Type = "A"
	req.Value = ip
	req.TTL = "600"

	_, err = client.UpdateDomainRecord(req)
	if err != nil {
		return err
	}

	return nil
}

// DoAlidnsV4 ...
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

	err = registerToAlidns(domain, subdomain, rid, ipv4)
	if err != nil {
		fmt.Println("register to f3322 error: ", err)
		return
	}

	// 重写ip缓存文件.
	dnsutils.FileWriteString("ipv4address", ipv4)
}
