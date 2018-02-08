package qiniustg

import (
	"config"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"strings"
)

func KeyToLower(recordsCh, retCh chan string, cfg config.Config) {

	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	stgCfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &stgCfg)

	for key := range recordsCh {

		lowCaseKey := strings.ToLower(key)
		if key == lowCaseKey {
			retCh <- fmt.Sprintf("%s\t%s", key, "CASE_PASS")
			continue
		}

		err := bucketManager.Copy(cfg.Bucket, key, cfg.Bucket, lowCaseKey, false)
		retCh <- fmt.Sprintf("%s\t%s\t%v", key, "CASE_TO_LOW", err)
	}
}
