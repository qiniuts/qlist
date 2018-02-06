package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpReq(retCh chan string, urls []string, cfg config.Config) {

	for _, url := range urls {
		resp, err := http.Get(url + cfg.FopQuery)
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
	// return nil
}
