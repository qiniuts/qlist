package utils

import (
	"io"
	"math/rand"
	"net"
	"net/http"
	"time"
)

var dIps map[string][]string

func init() {
	dIps = map[string][]string{}
}

func DmIp(host string) (ip string, err error) {

	ips, ok := dIps[host]
	if !ok {
		ipEntries, err := net.LookupIP(host)
		if err != nil {
			return "", err
		}

		for _, ip := range ipEntries {
			ips = append(ips, ip.String())
		}

		dIps[host] = ips
	}

	rand.Seed(time.Now().UTC().UnixNano())

	return ips[rand.Intn(len(ips))], nil
}

func NewRequest(method, url1 string, body io.Reader, headerHost string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url1, body)
	if err != nil {
		return
	}

	ip, err := DmIp(req.Host)
	if err != nil {
		return
	}

	req.URL.Host = ip
	if headerHost != "" {
		req.Host = headerHost
	}

	return
}
