package qiniustg

import (
	"encoding/base64"
	"encoding/json"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"log"
	"sync"
	"utils"
)

func (c *Client) List(inCh chan string) {

	mac := qbox.NewMac(c.AccessKey, c.SecretKey)
	bucketMgr := storage.NewBucketManager(mac, nil)

	wg := sync.WaitGroup{}
	marker, _ := newestListMarker(c.ProcResultsPath)

	for {
		items, _, markerOut, hasNext, err := bucketMgr.ListFiles(c.Bucket, c.Prefix, "", marker, 1000)
		if err != nil {
			log.Println("ListFiles", err)
			continue
		}
		marker = markerOut

		wg.Add(1)
		go func(wgp *sync.WaitGroup, listItems []storage.ListItem) {
			defer wgp.Done()
			for _, item := range listItems {
				inCh <- item.Key
			}
		}(&wg, items)

		if !hasNext {
			break
		}
	}

	wg.Wait()
	close(inCh)
}

func newestListMarker(fpath string) (marker string, err error) {

	line, err := utils.FileLastLine(fpath)
	if err != nil {
		return
	}

	mk := struct {
		C int    `json:"c"`
		K string `json:"k"`
	}{}

	mk.K = string(line)

	jmk, err := json.Marshal(mk)
	if err != nil {
		return
	}
	marker = base64.URLEncoding.EncodeToString(jmk)

	return
}
