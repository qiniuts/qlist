package qiniustg

import (
	"config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Hash struct {
	Md5 string `json:"md5"`
}

func Md5(recordsCh, retCh chan string, cfg config.Config) {

	hash := Hash{}

	for item := range recordsCh {

		fs := strings.Fields(item)
		if len(fs) < 4 {
			retCh <- fmt.Sprintf("%d\t%s\t%s", 900, item, "invalid item")
			continue
		}

		key := fs[0]
		createT := fs[3]
		url1 := "http://img.momocdn.com/" + key

		resp, err := http.Get(url1 + cfg.FopQuery)
		if err != nil {
			retCh <- fmt.Sprintf("%d\t%s\t%s", 900, key, err.Error())
			continue
		}

		err = json.NewDecoder(resp.Body).Decode(&hash)
		resp.Body.Close()
		if err != nil {
			retCh <- fmt.Sprintf("%d\t%s\t%s", 900, key, err.Error())
			continue
		}

		retCh <- fmt.Sprintf("%d\t%s\t%s\t%s", resp.StatusCode, key, createT, hash.Md5)
	}
}
