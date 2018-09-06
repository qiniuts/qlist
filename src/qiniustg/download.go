package qiniustg

import (
	"config"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Download(recordsCh, retCh chan string, cfg config.Config) {

	for url := range recordsCh {

		req, err := http.NewRequest("GET", url+cfg.FopQuery, nil)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		defer resp.Body.Close()
		code := resp.StatusCode

		//str := strings.Split(resp.Request.URL.Path, "/")
		fmt.Println("req.URL.Path:", req.URL.Path)

		fkey := base64.StdEncoding.EncodeToString([]byte(req.URL.Path))
		f, err := os.Create(fkey)

		if err != nil {
			panic(err)
		}

		fsize, err := io.Copy(f, resp.Body)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		if f.Close() != nil {
			retCh <- fmt.Sprintf("%s\t%d\t%s", url, 900, err.Error())
			continue
		}

		end := time.Now()
		sd := fmt.Sprintf("%d\t%d\t%d\t%.2f\n", start.Unix(), end.Unix(), fsize, float64(fsize)/end.Sub(start).Seconds())

		retCh <- fmt.Sprintf("%s\t%d\t%s", url, code, sd)
	}
}
