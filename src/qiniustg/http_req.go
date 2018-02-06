package qiniustg

import (
	"config"
	"fmt"
	"net/http"
)

func HttpReq(retCh chan string, urls []string, cfg config.Config) {

	for _, u := range urls {
		code, body, err := http.Get(nil, u+cfg.FopQuery)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", u, code, err.Error())
		} else {
			retCh <- fmt.Sprintf("%s\t%d\t%s", u, code, body)
		}
	}
}
