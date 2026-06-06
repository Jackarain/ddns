package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Jackarain/ddns/alidns"
	"github.com/Jackarain/ddns/cloudflare"
	"github.com/Jackarain/ddns/dnspod"
	"github.com/Jackarain/ddns/dnsutils"
	"github.com/Jackarain/ddns/f3322"
	"github.com/Jackarain/ddns/freedns"
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
	useFreeDNS    bool

	token  string
	user   string
	passwd string

	domain    string
	subdomain string
	dnsType   string

	command    string
	configFile string
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
	flag.BoolVar(&useFreeDNS, "freedns", false, "Use freedns api")

	flag.StringVar(&dnsutils.FetchIPv4AddrUrl, "externalIPv4", "", "Provide a URL to get the external IPv4 address")
	flag.StringVar(&dnsutils.FetchIPv6AddrUrl, "externalIPv6", "", "Provide a URL to get the external IPv6 address")

	// token 用于 dnspod, godaddy, namesilo, cloudflare, henet, freedns api
	flag.StringVar(&token, "token", "", "godaddy api-key:secret-key, dnspod token, namesilo api-key:secret-key, cloudflare api-token, henet password, freedns token")

	// user, passwd 用于 f3322/oray/ali api
	flag.StringVar(&user, "user", "", "f3322/oray/ali username only")
	flag.StringVar(&passwd, "passwd", "", "f3322/oray/ali password only")

	flag.StringVar(&domain, "domain", "", "Main domain")
	flag.StringVar(&subdomain, "subdomain", "", "Sub domain")
	flag.StringVar(&dnsType, "dnstype", "A", "dns type, AAAA/A")

	flag.StringVar(&command, "command", "", "Use command's output as IP address")

	flag.StringVar(&configFile, "config", "", "Path to config file")
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
	switch dnsType {
	case "A":
		dnspod.DoDNSPODv4(domain, subdomain, token, rid, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		godaddy.DoGodaddyv4(domain, subdomain, token, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		if len(subdomain) > 0 && len(domain) > 0 {
			domain = subdomain + "." + domain
		}
		f3322.DoF3322v4(domain, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		namesilo.DoNamesiloV4(domain, subdomain, token, rid, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		henet.DoHenetv4(domain, subdomain, token, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		if len(subdomain) > 0 && len(domain) > 0 {
			domain = subdomain + "." + domain
		}
		oray.DoOrayv4(domain, extIP)
	case "AAAA":
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

	switch dnsType {
	case "A":
		alidns.DoAlidnsV4(domain, subdomain, rid, extIP)
	case "AAAA":
		alidns.DoAlidnsV6(domain, subdomain, rid, extIP)
	}
}

func doFreeDNS() {
	if len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("freedns token/dnstype required")
		return
	}

	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	switch dnsType {
	case "A":
		freedns.DoFreeDNSv4(token, extIP)
	case "AAAA":
		freedns.DoFreeDNSv6(token, extIP)
	}
}

func doCloudFlare() {
	if len(domain) == 0 || len(subdomain) == 0 || len(token) == 0 || len(dnsType) == 0 {
		fmt.Println("cloudflare domain/subdomain/token/dnstype required")
		return
	}

	// 首先从文件中读取 zone_id, 并缓存到本地文件.
	zoneIDFileName := subdomain + "zone_id"
	zone_id, err := dnsutils.FileReadString(zoneIDFileName)
	if err != nil || zone_id == "" {
		zone_id, err = cloudflare.FetchZoneID(domain, token)
		if err != nil {
			fmt.Println("FetchZoneID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(zoneIDFileName, zone_id)
	}

	// 如果指定了command, 则使用command的输出内容作为公网ip
	var extIP string
	if command != "" {
		extIP = dnsutils.DoCommand(command)
		fmt.Println(extIP)
	}

	if len(subdomain) > 0 && len(domain) > 0 {
		domain = subdomain + "." + domain
	}

	// 从文件中读取record id
	ridFileName := subdomain + dnsType
	rid, err := dnsutils.FileReadString(ridFileName)
	if err != nil || rid == "" {
		rid, err = cloudflare.FetchRecordID(zone_id, token, domain)
		if err != nil {
			fmt.Println("FetchRecordID: " + err.Error())
			return
		}

		dnsutils.FileWriteString(ridFileName, rid)
	}

	// 根据dnsType选择更新A记录或AAAA记录
	switch dnsType {
	case "A":
		cloudflare.DoCFv4(domain, token, zone_id, rid, extIP)
	case "AAAA":
		cloudflare.DoCFv6(domain, token, zone_id, rid, extIP)
	}
}

// findConfig 从 os.Args 中查找 -config 参数，返回配置文件路径
// 并从 os.Args 中移除 -config 及其值，避免影响后续 flag.Parse()
func findConfig() string {
	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-config" || arg == "--config" {
			if i+1 < len(os.Args) {
				cfg := os.Args[i+1]
				os.Args = append(os.Args[:i], os.Args[i+2:]...)
				return cfg
			}
		}
		if strings.HasPrefix(arg, "-config=") || strings.HasPrefix(arg, "--config=") {
			cfg := strings.SplitN(arg, "=", 2)[1]
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			return cfg
		}
	}
	return ""
}

// loadConfig 从配置文件中读取参数并设置对应的变量。
// 配置文件每行一个参数，支持 # 开头的注释行，
// 格式为 key=value（如 domain=example.com）或 key（布尔值设为 true）。
// 参数名称与命令行参数名称一致（不需要前面的 - 前缀）。
func loadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var key, value string
		if idx := strings.Index(line, "="); idx != -1 {
			key = strings.TrimSpace(line[:idx])
			value = strings.TrimSpace(line[idx+1:])
		} else {
			key = line
			value = "true"
		}

		switch key {
		case "godaddy":
			useGodaddy = value == "true"
		case "dnspod":
			useDnspod = value == "true"
		case "f3322":
			useF3322 = value == "true"
		case "oray":
			useOray = value == "true"
		case "namesilo":
			useNamesilo = value == "true"
		case "henet":
			useHenet = value == "true"
		case "ali":
			useAlidns = value == "true"
		case "cloudflare":
			useCloudFlare = value == "true"
		case "freedns":
			useFreeDNS = value == "true"
		case "externalIPv4":
			dnsutils.FetchIPv4AddrUrl = value
		case "externalIPv6":
			dnsutils.FetchIPv6AddrUrl = value
		case "token":
			token = value
		case "user":
			user = value
		case "passwd":
			passwd = value
		case "domain":
			domain = value
		case "subdomain":
			subdomain = value
		case "dnstype":
			dnsType = value
		case "command":
			command = value
		default:
			fmt.Printf("config: unknown key %q at line %d, ignored\n", key, lineNo)
		}
	}
	return scanner.Err()
}

func main() {
	// 先查找并加载配置文件（如果指定了 -config）
	if cfgPath := findConfig(); cfgPath != "" {
		if err := loadConfig(cfgPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", cfgPath, err)
			os.Exit(1)
		}
	}

	flag.Parse()
	if help || len(os.Args) == 1 {
		// 如果是通过配置文件设置了参数，len(os.Args)==1 时仍然应该继续执行
		if configFile == "" && !useGodaddy && !useDnspod && !useF3322 && !useOray &&
			!useNamesilo && !useHenet && !useAlidns && !useCloudFlare && !useFreeDNS {
			flag.Usage()
			return
		}
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
	} else if useFreeDNS { // freedns api
		doFreeDNS()
	} else {
		fmt.Println("No api selected")
	}
}
