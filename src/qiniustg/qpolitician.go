package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Qpolitician(retCh chan string, keys []string, _ config.Config) error {

	for _, key := range keys {
		resp, err := http.Get(key + "?qpolitician")
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", key, 900, err.Error())
			continue
		}

		defer resp.Body.Close()
		code := resp.StatusCode
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", key, 900, err.Error())
			continue
		}

		retCh <- fmt.Sprintf("%s\t%d\t%s", key, code, body)
	}
	return nil
}
