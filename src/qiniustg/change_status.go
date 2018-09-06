package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/storage"
)

func ChangeFileStatus(retCh chan string, keys []string, cfg config.Config) {

	bucketManager := NewQNClient(cfg).BucketMgr()
	chstatusOps := []string{}

	for _, key := range keys {
		chstatusOps = append(chstatusOps, URIChangeStatus(cfg.Bucket, key, 1))
	}

	rets, err := bucketManager.Batch(chstatusOps)
	if err != nil {
		fmt.Printf("Batch Error: %#v", err)
	}
	for i, ret := range rets {
		retCh <- fmt.Sprintf("%d\t%s\t%s", ret.Code, keys[i], ret.Data.Error)
	}
}

func URIChangeStatus(bucket, key string, fileType int) string {
	return fmt.Sprintf("/chstatus/%s/status/%d", storage.EncodedEntry(bucket, key), fileType)
}
