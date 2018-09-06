package qiniustg

import (
	"config"
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io/ioutil"
)

func Adl(recordsCh, retCh chan string, cfg config.Config) {

	cli := NewAQNClient(cfg)
	uid_, ok := cfg.GetExtraParam("adl", "uid")
	fmt.Println(cfg)
	if !ok {
		panic("no uid in config")
	}
	uid := uid_.(string)

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

type AQNClient struct {
	*qbox.Mac
	*storage.Client
}

func NewAQNClient(cfg config.Config) *AQNClient {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	return &AQNClient{mac, &storage.DefaultClient}
}

func (cli *AQNClient) get(uid, bucket, key string) ([]byte, error) {

	url1 := "https://iovip.qbox.me/adminget/"
	params := map[string][]string{
		"uid":    {uid},
		"bucket": {bucket},
		"key":    {key},
	}
	ctx := context.WithValue(context.TODO(), "mac", cli.Mac)

	resp, err := cli.DoRequestWithForm(ctx, "POST", url1, nil, params)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
