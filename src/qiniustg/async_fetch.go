package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/storage"
	"strings"
)

func AsyncFetch(recordsCh, retCh chan string, cfg config.Config) {

	bucketManager := NewClient(cfg).BucketMgr()
	param := storage.AsyncFetchParam{}

	for record := range recordsCh {

		fs := strings.Split(record, "\t")
		param.Url = fs[0]
		param.Bucket = cfg.Bucket

		if len(fs) > 1 {
			param.Key = fs[1]
		}

		if len(fs) > 2 {
			param.Md5 = fs[2]
		}

		ret, err := bucketManager.AsyncFetch(param)
		retCh <- fmt.Sprintf("%v\t%v", ret, err)
	}
}
