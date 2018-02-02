package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func qpulp(retCh chan string, keys []string, _ config.Config) error {

	for _, key := range keys {
		resp, _ := http.Get(str + "?qpulp")
		key := str
		code := resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)
		retCh <- fmt.Sprintf("%s\t%d\t%s\n", key, code, body)
		resp.Body.Close()
	}
	return err

}
