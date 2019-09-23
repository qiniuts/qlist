package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/storage"
)


func ChangeFileType(retCh chan string, keys []string, cfg config.Config) {

	rets, err := changeFileType(keys, cfg)
	if err != nil {
		fmt.Printf("Batch Error: %#v", err)
	}

	for i, ret := range rets {
		retCh <- fmt.Sprintf("%d\t%s\t%s", ret.Code, keys[i], ret.Data.Error)
	}
}

func changeFileType(keys []string, cfg config.Config) (rets []storage.BatchOpRet, err error) {

	bucketManager := NewQNClient(cfg).BucketMgr()

	chtypeOps := []string{}

	for _, key := range keys {
		chtypeOps = append(chtypeOps, storage.URIChangeType(cfg.Bucket, key, 1))
	}

	rets, err = bucketManager.Batch(chtypeOps)
	return
}
