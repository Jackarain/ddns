package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jackarain/ddns/dnspod"
	"github.com/Jackarain/ddns/dnsutils"
	"github.com/Jackarain/ddns/f3322"
	"github.com/Jackarain/ddns/godaddy"
	"github.com/Jackarain/ddns/henet"
	"github.com/Jackarain/ddns/namesilo"
)

var (
	help bool

	useGodaddy  bool
	useDnspod   bool
	useF3322    bool
	useNamesilo bool
	useHenet    bool

	token     string
	domain    string
	subdomain string
	dnsType   string

	command string

	extIP string
)

func init() {
	flag.BoolVar(&help, "help", false, "help message")

	flag.BoolVar(&useGodaddy, "godaddy", false, "Use godaddy api")
	flag.BoolVar(&useDnspod, "dnspod", false, "Use dnspod api")
	flag.BoolVar(&useF3322, "f3322", false, "Use f3322 api")
	flag.BoolVar(&useNamesilo, "namesilo", false, "Use namesilo api")
	flag.BoolVar(&useHenet, "henet", false, "Use henet api")

	flag.StringVar(&dnsutils.FetchIPv4AddrUrl, "externalIPv4", "", "Provide a URL to get the external IPv4 address")
	flag.StringVar(&dnsutils.FetchIPv6AddrUrl, "externalIPv6", "", "Provide a URL to get the external IPv6 address")

	// token 用于 dnspod, godaddy, namesilo, henet api
	flag.StringVar(&token, "token", "", "godaddy api-key:secret-key, dnspod token, namesilo api-key:secret-key, henet password")
	// f3322user, f3322passwd 用于 f3322 api
	flag.StringVar(&f3322.User, "f3322user", "", "f3322 username")
	flag.StringVar(&f3322.Passwd, "f3322passwd", "", "f3322 password")

	flag.StringVar(&domain, "domain", "", "Main domain")
	flag.StringVar(&subdomain, "subdomain", "", "Sub domain")
	flag.StringVar(&dnsType, "dnstype", "A", "dns type, AAAA/A")

	flag.StringVar(&command, "command", "", "ip use command result")
}

func doDnspod() {
	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = dnspod.FetchRecordID(token, domain, subdomain, dnsType)
		if err != nil {
			fmt.Println(err)
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
	if len(f3322.User) == 0 || len(f3322.Passwd) == 0 {
		fmt.Println("f3322user and f3322passwd required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	if dnsType == "A" {
		f3322.DoF3322v4(domain, extIP)
	} else if dnsType == "AAAA" {
		fmt.Println("f3322 doesn’t work with ipv6")
	}
}

func doNamesilo() {
	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = namesilo.FetchRecordID(token, domain, subdomain)
		if err != nil {
			fmt.Println(err)
			return
		}

		// dnsutils.FileWriteString(ridFileName, rid)
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
	} else if useNamesilo { // namesilo api
		doNamesilo()
	} else if useHenet { // henet api
		doHenet()
	} else {
		fmt.Println("No api selected")
	}
}
