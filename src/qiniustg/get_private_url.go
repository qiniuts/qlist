package qiniustg

import (
	"config"
	"fmt"
	"strings"
	"github.com/qiniu/api.v7/storage"
	"time"
)

func GetPrivateUrl(recordsCh, retCh chan string, cfg config.Config) {

	mac := NewQNClient(cfg).StorgeMgr()

	//mac := auth.New("m7Vt8XvmWXijN_szFoFQHywTkuSciG2CdpZ3pIX8", "")

	deadline := time.Now().Add(time.Second * 3600 * 1000).Unix() //1小时有效期

	for record := range recordsCh {

		fs := strings.Split(record, "\t")

		fmt.Println(record)

		privateAccessURL := storage.MakePrivateURL(mac, "domain", fs[0], deadline)
		if len(fs) >= 3 {
			retCh <- fmt.Sprintf("%s\t%s\t%s", privateAccessURL, fs[0],fs[2])
		}
	}

}