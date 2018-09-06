package utils

import (
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

var dIps sync.Map

func DmIp(host string) (ip string, err error) {

	ips, ok := dIps.Load(host)
	if !ok {

		ipEntries, err := net.LookupIP(host)
		if err != nil {
			return "", err
		}
		for _, ip := range ipEntries {
			ips = append(ips.([]string), ip.String())
		}

		dIps.Store(host, ips)
	}

	dmIps := []string(ips.([]string))
	rand.Seed(time.Now().UTC().UnixNano())

	return dmIps[rand.Intn(len(dmIps))], nil
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
