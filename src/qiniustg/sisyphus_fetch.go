package qiniustg

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"config"
	"net/url"
	"strconv"
)

type AsyncFetchRet struct {
	Id string `json:"id"`
	Wait int `json:"wait"`
}

type AsyncFetchParam struct {
	Url string `json:"url"`
	Host string `json:"host"`
	Bucket string `json:"bucket"`
	Key string `json:"key"`
}


func AsyncFetch(failedCh chan string, keys []string, cfg config.Config) error  {

	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	cli := storage.NewClient(mac, nil)
	u := "http://api-z2.qiniu.com/sisyphus/fetch"

	for _, key := range keys {

		retStr := ""
		u_, err := url.Parse(key)
		if err != nil {
			retStr = err.Error() + "\t" + key
			failedCh <- retStr
			continue
		}

		ret := AsyncFetchRet{}
		param := AsyncFetchParam{Url:key, Bucket: cfg.Bucket, Key: u_.Path[1:]}

		err = cli.CallWithJson(nil, &ret, "POST", u, param)
		if err != nil {
			retStr = err.Error()
		} else {
			retStr = ret.Id + "\t" + strconv.Itoa(ret.Wait)
		}

		retStr += "\t" + key
		failedCh <- retStr
	}

	return nil
}


