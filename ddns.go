package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Jackarain/ddns/dnspod"
	"github.com/Jackarain/ddns/dnsutils"
	"github.com/Jackarain/ddns/godaddy"
)

var (
	help       bool
	useGodaddy bool
	useDnspod  bool
	token      string
	domain     string
	subdomain  string
	dnsType    string

	command string
	args    string

	extIP string
)

func init() {
	flag.BoolVar(&help, "help", false, "help message")
	flag.BoolVar(&useGodaddy, "godaddy", false, "Use godaddy api")
	flag.BoolVar(&useDnspod, "dnspod", false, "Use dnspod api")
	flag.StringVar(&token, "token", "", "Api token/secret,godaddy api-key:secret")
	flag.StringVar(&domain, "domain", "", "Main domain")
	flag.StringVar(&subdomain, "subdomain", "", "Sub domain")
	flag.StringVar(&dnsType, "dnstype", "", "dns type, AAAA/A")

	flag.StringVar(&command, "command", "", "ip use command result")
	flag.StringVar(&args, "args", "", "command args")
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
		extIP = dnsutils.DoCommand(command, args)
	}

	fmt.Println(subdomain, "dnspod record id:", rid)
	if dnsType == "A" {
		dnspod.DoDNSPODv4(domain, subdomain, rid, extIP)
	} else if dnsType == "AAAA" {
		dnspod.DoDNSPODv6(domain, subdomain, rid, extIP)
	}
}

func doGodaddy() {
	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command, args)
		fmt.Println(extIP)
	}

	if dnsType == "A" {
		godaddy.DoGodaddyv4(domain, subdomain, token, extIP)
	} else if dnsType == "AAAA" {
		godaddy.DoGodaddyv6(domain, subdomain, token, extIP)
	}
}

func main() {
	flag.Parse()
	if help || len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if useDnspod {
		doDnspod()
	} else if useGodaddy {
		doGodaddy()
	}
}
