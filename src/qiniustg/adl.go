package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"io/ioutil"
)

func Adl(recordsCh, retCh chan string, cfg config.Config) {

	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	cli := AClient{storage.NewClient(mac, nil)}
	uid := "1380703881"

	for key := range recordsCh {

		body, err := cli.get(uid, cfg.Bucket, key)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%s\t%v", key, "GET_ERR:", err)
			continue
		}

		err = ioutil.WriteFile(key, body, 0644)
		if err != nil {
			retCh <- fmt.Sprintf("%s\t%s\t%v", key, "WRITEFILE_ERR:", err)
			continue
		}
	}
}

type AClient struct {
	*rpc.Client
}

func (cli *AClient) get(uid, bucket, key string) ([]byte, error) {

	url1 := "https://iovip.qbox.me/aget/"
	params := map[string][]string{
		"uid":    {uid},
		"bucket": {bucket},
		"key":    {key},
	}
	resp, err := cli.DoRequestWithForm(nil, "POST", url1, params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
