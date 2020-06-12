package ddnsreg

import (
	"net/http"
	"net/url"

	"rpi4.p2sp.net/Jackarain/externalip"
)

// RegisterToF3322
func RegisterToF3322(ip string) error {
	f3322Url := "http://members.3322.net/dyndns/update?system=dyndns&hostname=sgrc.f3322.net&myip=" + ip
	request, err := http.NewRequest("GET", f3322Url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Basic "+externalip.BasicAuth("wgm001", "ggfggc"))
	f3322Client := &http.Client{}
	response, err := f3322Client.Do(request)
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

// IPv6RegisterToDNSPOD
func IPv6RegisterToDNSPOD(domain, subdomain, rid, ip string) error {
	dnspodURL := "https://dnsapi.cn/Record.Modify"
	response, err := http.PostForm(dnspodURL, url.Values{
		"login_token": {"11898,8d7347cd5969f7aa89752c068a6b949a"},
		"format":      {"json"},
		"domain":      {domain},
		"record_id":   {rid},
		"sub_domain":  {subdomain},
		"record_type": {"AAAA"},
		"record_line": {"默认"},
		"value":       {ip},
	})
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

// IPv4RegisterToDNSPOD
func IPv4RegisterToDNSPOD(domain, subdomain, rid, ip string) error {
	dnspodURL := "https://dnsapi.cn/Record.Modify"
	response, err := http.PostForm(dnspodURL, url.Values{
		"login_token": {"11898,8d7347cd5969f7aa89752c068a6b949a"},
		"format":      {"json"},
		"domain":      {domain},
		"record_id":   {rid},
		"sub_domain":  {subdomain},
		"record_type": {"A"},
		"record_line": {"默认"},
		"value":       {ip},
	})
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

