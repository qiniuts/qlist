package qiniustg

import (
	"config"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"utils"
)

func Download(recordsCh, retCh chan string, cfg config.Config) {

	for url := range recordsCh {

		req, err := utils.NewRequest("GET", url+cfg.FopQuery, nil, cfg.ReqHeaderHost)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		defer resp.Body.Close()
		code := resp.StatusCode

		str := strings.Split(resp.Request.URL.Path, "/")
		f, err := os.Create(str[1])

		if err != nil {
			panic(err)
		}

		io.Copy(f, resp.Body)
		f.Close()
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		retCh <- fmt.Sprintf("%s\t%d", url, code)
	}
}
