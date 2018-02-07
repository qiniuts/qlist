package localstg

import (
	"config"
)

func BucketList(recordsCh, retCh chan string, _ config.Config) {
	for record := range recordsCh {
		retCh <- record
	}
}
