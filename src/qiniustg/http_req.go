package qiniustg

import (
	"config"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpReq(recordsCh, retCh chan string, cfg config.Config) {

	for url := range recordsCh {

		req, err := http.NewRequest("GET", url+cfg.FopQuery, nil)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		retCh <- fmt.Sprintf("%s\t%d\t%s", url, resp.StatusCode, body)
	}
}
