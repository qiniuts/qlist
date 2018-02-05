package qiniustg

import (
	"config"
	"fmt"
	"github.com/valyala/fasthttp"
)

func HttpReq(retCh chan string, urls []string, cfg config.Config) error {

	for _, u := range urls {
		code, body, err := fasthttp.Get(nil, u+cfg.FopQuery)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", u, code, err.Error())
		} else {
			retCh <- fmt.Sprintf("%s\t%d\t%s", u, code, body)
		}
	}
	return nil
}
