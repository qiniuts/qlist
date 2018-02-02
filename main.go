package main

import (
	"config"
	"flag"
	"log"
	"qiniustg"
	"sync"
	"utils"
)

type batchProcFunc func(failedCh chan string, keys []string, cfg config.Config) error

func main() {

	cfgPath := flag.String("cfg_path", "cfg.json", "")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	cli := Client{cfg}
	qiniuCli := qiniustg.NewClient(cfg)

	inCh := make(chan string, 1000*100)
	doneCh := make(chan string, 1000*10)
	failedCh := make(chan string, 1000*10)

	go qiniuCli.List(inCh)
	go cli.Proc(inCh, doneCh, failedCh, qiniustg.qpulp)

	doneLogWait := make(chan bool)
	failedLogWait := make(chan bool)

	go utils.Log(doneCh, cli.DoneRecordPath, doneLogWait)
	go utils.Log(failedCh, "errors.log", failedLogWait)

	<-doneLogWait
	<-failedLogWait
	log.Println("list and proc done!")
}

type Client struct {
	config.Config
}

func (c *Client) Proc(inCh, doneCh, failedCh chan string, batch batchProcFunc) {
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

func (c *Client) worker(keysCh, processedCh, failedCh chan string, batch batchProcFunc, wg *sync.WaitGroup) {

	defer wg.Done()
	keys := []string{}
	for {
		key, ok := <-keysCh
		if ok {
			processedCh <- key
			keys = append(keys, key)
			if len(keys) < 1000 {
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
