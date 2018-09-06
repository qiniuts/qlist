package qiniustg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/log.v7"
	"strings"
	"sync"
	"utils"
)

func (c *QNClient) List(inCh chan string) {

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

func (c *QNClient) List2(inCh chan string) {

	mac := qbox.NewMac(c.AccessKey, c.SecretKey)
	bucketMgr := storage.NewBucketManager(mac, nil)
	marker, _ := newestListMarker(c.ProcResultsPath)
	defer close(inCh)

	for {
		retChan, err := bucketMgr.ListBucket(c.Bucket, c.Prefix, "", marker)
		if err != nil {
			log.Error("ListFiles Error:", err, c.Bucket, marker)
			continue
		}

		for ret := range retChan {
			marker = ret.Marker
			if ret.Item.Key == "" {
				continue
			}

			inCh <- fmt.Sprintf("%s\t%d\t%d", ret.Item.Key, ret.Item.Fsize, ret.Item.PutTime)
			if marker == "" {
				return
			}
		}
	}
}

func newestListMarker(fpath string) (marker string, err error) {

	line, err := utils.FileLastLine(fpath)
	if err != nil || string(line) == "" {
		return
	}

	mk := struct {
		C int    `json:"c"`
		K string `json:"k"`
	}{}

	ls := strings.Fields(string(line))
	if len(ls) < 1 {
		return
	}
	mk.K = ls[0]

	jmk, err := json.Marshal(mk)
	if err != nil {
		return
	}
	marker = base64.URLEncoding.EncodeToString(jmk)

	return
}
