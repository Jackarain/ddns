package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jackarain/ddns/alidns"
	"github.com/Jackarain/ddns/cloudflare"
	"github.com/Jackarain/ddns/dnspod"
	"github.com/Jackarain/ddns/dnsutils"
	"github.com/Jackarain/ddns/f3322"
	"github.com/Jackarain/ddns/godaddy"
	"github.com/Jackarain/ddns/henet"
	"github.com/Jackarain/ddns/namesilo"
	"github.com/Jackarain/ddns/oray"
)

var (
	help bool

	useGodaddy    bool
	useDnspod     bool
	useF3322      bool
	useOray       bool
	useNamesilo   bool
	useHenet      bool
	useAlidns     bool
	useCloudFlare bool

	token  string
	user   string
	passwd string

	domain    string
	subdomain string
	dnsType   string

	command string
)

func init() {
	flag.BoolVar(&help, "help", false, "help message")

	flag.BoolVar(&useGodaddy, "godaddy", false, "Use godaddy api")
	flag.BoolVar(&useDnspod, "dnspod", false, "Use dnspod api")
	flag.BoolVar(&useF3322, "f3322", false, "Use f3322 api")
	flag.BoolVar(&useOray, "oray", false, "Use oray api")
	flag.BoolVar(&useNamesilo, "namesilo", false, "Use namesilo api")
	flag.BoolVar(&useHenet, "henet", false, "Use henet api")
	flag.BoolVar(&useAlidns, "ali", false, "Use alidns api")
	flag.BoolVar(&useCloudFlare, "cloudflare", false, "Use cloudflare api")

	flag.StringVar(&dnsutils.FetchIPv4AddrUrl, "externalIPv4", "", "Provide a URL to get the external IPv4 address")
	flag.StringVar(&dnsutils.FetchIPv6AddrUrl, "externalIPv6", "", "Provide a URL to get the external IPv6 address")

	// token 用于 dnspod, godaddy, namesilo, cloudflare, henet api
	flag.StringVar(&token, "token", "", "godaddy api-key:secret-key, dnspod token, namesilo api-key:secret-key, cloudflare zone_id:api-key, henet password")

	// user, passwd 用于 f3322/oray/ali api
	flag.StringVar(&user, "user", "", "f3322/oray/ali username only")
	flag.StringVar(&passwd, "passwd", "", "f3322/oray/ali password only")

	flag.StringVar(&domain, "domain", "", "Main domain")
	flag.StringVar(&subdomain, "subdomain", "", "Sub domain")
	flag.StringVar(&dnsType, "dnstype", "A", "dns type, AAAA/A")

	flag.StringVar(&command, "command", "", "Use command's output as IP address")
}

func doDnspod() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("dnspod domain/subdomain/token/dnstype required")
		return
	}

	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = dnspod.FetchRecordID(token, domain, subdomain, dnsType)
		if err != nil {
			fmt.Println("FetchRecordID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(ridFileName, rid)
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	fmt.Println(subdomain, "dnspod record id:", rid)
	if dnsType == "A" {
		dnspod.DoDNSPODv4(domain, subdomain, token, rid, extIP)
	} else if dnsType == "AAAA" {
		dnspod.DoDNSPODv6(domain, subdomain, token, rid, extIP)
	}
}

func doGodaddy() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("godaddy domain/subdomain/token/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	if dnsType == "A" {
		godaddy.DoGodaddyv4(domain, subdomain, token, extIP)
	} else if dnsType == "AAAA" {
		godaddy.DoGodaddyv6(domain, subdomain, token, extIP)
	}
}

func doF3322() {
	if len(domain) == 0 || len(user) == 0 || len(passwd) == 0 || len(dnsType) == 0 {
		fmt.Println("f3322 domain/user/passwd/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	f3322.User = user
	f3322.Passwd = passwd

	if dnsType == "A" {
		if len(subdomain) > 0 && len(domain) > 0 {
			domain = subdomain + "." + domain
		}

		f3322.DoF3322v4(domain, extIP)
	} else if dnsType == "AAAA" {
		fmt.Println("f3322 doesn’t work with ipv6")
	}
}

func doNamesilo() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("namesilo domain/subdomain/token/dnstype required")
		return
	}

	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = namesilo.FetchRecordID(token, domain, subdomain)
		if err != nil {
			fmt.Println("FetchRecordID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(ridFileName, rid)
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	fmt.Println(subdomain, "namesilo record id:", rid)

	if dnsType == "A" {
		namesilo.DoNamesiloV4(domain, subdomain, token, rid, extIP)
	} else if dnsType == "AAAA" {
		namesilo.DoNamesiloV6(domain, subdomain, token, rid, extIP)
	}
}

func doHenet() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("henet domain/subdomain/token/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	if dnsType == "A" {
		henet.DoHenetv4(domain, subdomain, token, extIP)
	} else if dnsType == "AAAA" {
		henet.DoHenetv6(domain, subdomain, token, extIP)
	}
}

func doOray() {
	if len(domain) == 0 || len(subdomain) == 0 || len(user) == 0 || len(passwd) == 0 || len(dnsType) == 0 {
		fmt.Println("oray domain/subdomain/user/passwd/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	oray.User = user
	oray.Passwd = passwd

	if dnsType == "A" {
		if len(subdomain) > 0 && len(domain) > 0 {
			domain = subdomain + "." + domain
		}
		oray.DoOrayv4(domain, extIP)
	} else if dnsType == "AAAA" {
		fmt.Println("oray doesn’t work with ipv6")
	}
}

func doAlidns() {
	if len(domain) == 0 || len(subdomain) == 0 || len(user) == 0 || len(passwd) == 0 || len(dnsType) == 0 {
		fmt.Println("alidns domain/subdomain/user/passwd/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	alidns.User = user
	alidns.Passwd = passwd

	// 从文件中读取record id
	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = alidns.FetchRecordID(domain)
		if err != nil {
			fmt.Println("FetchRecordID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(ridFileName, rid)
	}

	if dnsType == "A" {
		alidns.DoAlidnsV4(domain, subdomain, rid, extIP)
	} else if dnsType == "AAAA" {
		alidns.DoAlidnsV6(domain, subdomain, rid, extIP)
	}
}

func doCloudFlare() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("cloudflare domain/subdomain/token/dnstype required")
		return
	}

	// 如果指定了command, 则使用command的输出内容作为公网ip
	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	// 从token中获取zone_id和api-key，格式为zone_id:api-key
	zone_id, api_key := dnsutils.ParseToken(token)

	if len(subdomain) > 0 && len(domain) > 0 {
		domain = subdomain + "." + domain
	}

	// 从文件中读取record id
	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = cloudflare.FetchRecordID(zone_id, api_key, domain)
		if err != nil {
			fmt.Println("FetchRecordID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(ridFileName, rid)
	}

	// 根据dnsType选择更新A记录或AAAA记录
	if dnsType == "A" {
		cloudflare.DoCFv4(domain, api_key, zone_id, rid, extIP)
	} else if dnsType == "AAAA" {
		cloudflare.DoCFv6(domain, api_key, zone_id, rid, extIP)
	}
}

func main() {
	flag.Parse()
	if help || len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if useDnspod { // dnspod api
		doDnspod()
	} else if useGodaddy { // godaddy api
		doGodaddy()
	} else if useF3322 { // f3322 api
		doF3322()
	} else if useOray { // oray api
		doOray()
	} else if useNamesilo { // namesilo api
		doNamesilo()
	} else if useHenet { // henet api
		doHenet()
	} else if useAlidns { // alidns api
		doAlidns()
	} else if useCloudFlare { // cloudflare api
		doCloudFlare()
	} else {
		fmt.Println("No api selected")
	}
}
