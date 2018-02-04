package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Qpolitician(retCh chan string, keys []string, _ config.Config) error {

	for _, key := range keys {
		resp, _ := http.Get(key + "?qpolitician")
		defer resp.Body.Close()

		code := resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)

		retCh <- fmt.Sprintf("%s\t%d\t%s", key, code, body)
	}
	return nil
}
