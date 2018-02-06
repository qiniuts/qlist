package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpReq(retCh chan string, urls []string, cfg config.Config) {

	for _, key := range keys {
		resp, err := http.Get(key + cfg.FopQuery)
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
