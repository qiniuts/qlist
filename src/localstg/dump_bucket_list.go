package localstg

import (
	"config"
)

func DumpList(recordsCh, retCh chan string, _ config.Config) {
	for record := range recordsCh {
		retCh <- record
	}
}
