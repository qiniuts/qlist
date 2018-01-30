package main

import (
	"flag"
	"os"
	"bufio"
	"log"
	"sync"
	"qiniustg"
	"config"
	"localstg"
)

type batchProcFunc func(failedCh chan string, keys []string, cfg config.Config) error

func main()  {

	cfgPath := flag.String("cfg_path", "cfg.json", "")
	fileToProcPath := flag.String("file_path", "keys_path.txt", "")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	cli := Client{cfg}

	inCh := make(chan string, 1000*100)
	doneCh := make(chan string, 1000*10)
	failedCh := make(chan string, 1000*10)

	//go cli.List(inCh)
	go localstg.List(inCh, *fileToProcPath)

	//go cli.Proc(inCh, doneCh, failedCh, qiniustg.ChangeFileType2)
	go cli.Proc(inCh, doneCh, failedCh, qiniustg.AsyncFetch)

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


