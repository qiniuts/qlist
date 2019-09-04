package qiniustg

import (
	"config"
	"fmt"
	"net/http"
)

func GetHttpCode(recordsCh, retCh chan string, cfg config.Config) {

	for url1 := range recordsCh {
		req, err := http.NewRequest("HEAD", url1, nil)

		if err != nil {
			retCh <- fmt.Sprintf("%sResultQiniu%s", url1,err)
		}else{
			resp, err2 := http.DefaultClient.Do(req)

			if err != nil {
				retCh <- fmt.Sprintf("%sResultQiniu%s", url1,err2)
			}else {
				retCh <- fmt.Sprintf("%sResultQiniu:%d", url1,resp.StatusCode)
			}
		}
	}
}
