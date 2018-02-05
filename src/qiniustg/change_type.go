package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

func ChangeFileType(retCh chan string, keys []string, cfg config.Config) {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	stgCfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &stgCfg)
	chtypeOps := []string{}

	for _, key := range keys {
		chtypeOps = append(chtypeOps, storage.URIChangeType(cfg.Bucket, key, 1))
	}

	rets, err := bucketManager.Batch(chtypeOps)
	if err != nil {
		fmt.Printf("Batch Error: %#v", err)
	}
	for i, ret := range rets {
		retCh <- fmt.Sprintf("%d\t%s\t%s", ret.Code, keys[i], ret.Data.Error)
	}
}
