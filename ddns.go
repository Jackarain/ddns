package main

import (
	"flag"
	"fmt"
	"os"

	"git.superpool.io/Jackarain/ddns/dnspod"
	"git.superpool.io/Jackarain/ddns/externalip"
	"git.superpool.io/Jackarain/ddns/godaddy"
)

var (
	help       bool
	useGodaddy bool
	useDnspod  bool
	token      string
	domain     string
	subdomain  string
	dnsType    string
)

func init() {
	flag.BoolVar(&help, "help", false, "help message")
	flag.BoolVar(&useGodaddy, "godaddy", false, "Use godaddy api")
	flag.BoolVar(&useDnspod, "dnspod", false, "Use dnspod api")
	flag.StringVar(&token, "token", "", "Api token/secret")
	flag.StringVar(&domain, "domain", "", "Main domain")
	flag.StringVar(&subdomain, "subdomain", "", "Sub domain")
	flag.StringVar(&dnsType, "dnstype", "", "dns type, AAAA/A")
}

func doDnspod() {
	ridFileName := subdomain + dnsType
	rid, err := externalip.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = dnspod.FetchRecordID(token, domain, subdomain, dnsType)
		if err != nil {
			fmt.Println(err)
			return
		}

		externalip.FileWriteString(ridFileName, rid)
	}

	fmt.Println(subdomain, "dnspod record id:", rid)
	if dnsType == "A" {
		dnspod.DoDNSPODv4(domain, subdomain, rid)
	} else if dnsType == "AAAA" {
		dnspod.DoDNSPODv6(domain, subdomain, rid)
	}
}

func doGodaddy() {
	if dnsType == "A" {
		godaddy.DoGodaddyv4(domain, subdomain, token)
	} else if dnsType == "AAAA" {
		godaddy.DoGodaddyv6(domain, subdomain, token)
	}
}

func main() {
	flag.Parse()
	if help || len(os.Args) == 1 {
		flag.Usage()
	}

	if useDnspod {
		doDnspod()
	} else if useGodaddy {
		doGodaddy()
	}
}
