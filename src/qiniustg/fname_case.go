package qiniustg

import (
	"config"
	"fmt"
	"strings"
)

func KeyToLower(recordsCh, retCh chan string, cfg config.Config) {

	bucketManager := NewClient(cfg).BucketMgr()
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
