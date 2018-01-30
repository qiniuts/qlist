package main

import (
	"flag"
	"os"
	"encoding/json"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/api.v7/auth/qbox"
	"bufio"
	"encoding/base64"
	"log"
	"sync"
	"qiniustg"
	"config"
)

type batchProcFunc func(failedCh chan string, keys []string, cfg config.Config) error

func main()  {

	cfgPath := flag.String("cfg_path", "cfg.json", "")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	cli := Client{cfg}

	inCh := make(chan string, 1000*100)
	doneCh := make(chan string, 1000*10)
	failedCh := make(chan string, 1000*10)

	go cli.List(inCh)
	go cli.Proc(inCh, doneCh, failedCh, qiniustg.ChangeFileType2)

	doneLogWait := make(chan bool)
	failedLogWait := make(chan bool)

	go cli.Log(doneCh, cli.DoneRecordPath, doneLogWait)
	go cli.Log(failedCh, "errors.log", failedLogWait)

	<-doneLogWait
	<-failedLogWait
	log.Println("list and proc done!")
}

type Client struct {
	config.Config
}

func (c *Client) Proc(inCh, doneCh, failedCh chan string, batch batchProcFunc)  {
	wg := sync.WaitGroup{}
	wg.Add(c.WorkerCount)

	for i := 0; i < c.WorkerCount; i++ {
		go c.worker(inCh, doneCh, failedCh, batch, &wg)
	}

	wg.Wait()
	close(doneCh)
	close(failedCh)

	return
}

func (c *Client) worker(keysCh, processedCh, failedCh chan string, batch batchProcFunc, wg *sync.WaitGroup)  {

	defer wg.Done()
	keys := []string{}
	for {
		key, ok := <- keysCh
		if ok {
			processedCh <- key
			keys = append(keys, key)
			if  len(keys) < 1000 {
				continue
			}
		} else if len(keys) == 0 {
			break
		}

		err := batch(failedCh, keys, c.Config)
		if err != nil {
			log.Println("Error", err)
		}

		keys = []string{}
	}
}


func (c *Client) List(inCh chan string)  {

	mac := qbox.NewMac(c.AccessKey, c.SecretKey)
	bucketMgr := storage.NewBucketManager(mac,nil)

	wg := sync.WaitGroup{}
	marker, _ := getNewestMarker(c.DoneRecordPath)

	for {
		items, _, markerOut, hasNext, err := bucketMgr.ListFiles(c.Bucket, "", "", marker, 1000)
		if err != nil {
			log.Println("ListFiles", err)
			continue
		}
		marker = markerOut

		wg.Add(1)
		go func(wgp *sync.WaitGroup, listItems []storage.ListItem) {
			defer wgp.Done()
			for _, item := range listItems {
				if item.Type != 1 {
					inCh <- item.Key
				}
			}
		}(&wg, items)

		if !hasNext {
			break
		}
	}

	wg.Wait()
	close(inCh)
}

func (c *Client)Log(outCh chan string, fpath string, done chan bool)  {
	fh, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	for key := range outCh {
		fh.WriteString(key + "\n")
	}
	close(done)
}

func getNewestMarker(fpath string) (marker string, err error) {

	line, err := tailLastLine(fpath)
	if err != nil {
		return
	}

	mk := struct {
		C int `json:"c"`
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

func tailLastLine(fpath string) (ll []byte, err error) {
	fi, err := os.Stat(fpath)
	if err != nil {
		return
	}
	fsize := fi.Size()

	fh, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer fh.Close()

	offset := fsize - 1
	buf := make([]byte, 1)
	for offset >= 0  {

		_, err = fh.ReadAt(buf, offset)
		if err != nil {
			return
		}

		if offset == 0 {
			ll, _, err = bufio.NewReader(fh).ReadLine()
			return
		}

		if string(buf) == "\n" && offset != fsize - 1{

			l := fsize - offset - 1
			lbuf := make([]byte, l)

			fh.ReadAt(lbuf, offset + 1)
			ll = lbuf
			return
		}
		offset--
	}

	return
}
