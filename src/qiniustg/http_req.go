package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func HttpReq(recordsCh, retCh chan string, cfg config.Config) {

	dIps := map[string][]string{}
	getIp := func(host string) (ip string, err error) {

		ips, ok := dIps[host]
		if !ok {
			ipEntries, err := net.LookupIP(host)
			if err != nil {
				return "", err
			}
			fmt.Println(ipEntries)

			for _, ip := range ipEntries {
				ips = append(ips, ip.String())
			}

			dIps[host] = ips
		}

		rand.Seed(time.Now().UTC().UnixNano())

		fmt.Println(ips)

		return ips[rand.Intn(len(ips))], nil
	}
	cli := http.DefaultClient

	for url := range recordsCh {

		req, err := http.NewRequest("GET", url+cfg.FopQuery, nil)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		ip, err := getIp(req.Host)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		fmt.Println(ip)
		req.URL.Host = ip

		resp, err := cli.Do(req)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		defer resp.Body.Close()
		code := resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		retCh <- fmt.Sprintf("%s\t%d\t%s", url, code, body)
	}
}
