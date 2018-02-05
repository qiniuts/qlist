package main

import (
	"config"
	"flag"
	"fmt"
	"localstg"
	"log"
	"qiniustg"
	"sync"
	"utils"
)

type batchProcFunc func(failedCh chan string, keys []string, cfg config.Config) error

var procFuncs = map[string]batchProcFunc{
	"req":      qiniustg.HttpReq,
	"chtype":   qiniustg.ChangeFileType,
	"chstatus": qiniustg.ChangeFileStatus,
}

func main() {

	cfgPath := flag.String("cfg_path", "cfg.json", "")
	filePath := flag.String("file_path", "keys.txt", "")

	flag.Parse()
	funcName := flag.Arg(0)

	procFunc, ok := procFuncs[funcName]
	if !ok {
		fmt.Printf("No %s Function", funcName)
		return
	}

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}
	cli := Client{cfg}

	//channels to cache records and proc result
	recordsCh := make(chan string, 1000*100)
	doneRecordsCh := make(chan string, 1000*10)
	procResultCh := make(chan string, 1000*10)

	//list records local file
	go localstg.List(recordsCh, *filePath)

	//list records qiniu storage
	//qiniuCli := qiniustg.NewClient(cfg)
	//go qiniuCli.List(recordsCh)

	//proc records
	go cli.Proc(recordsCh, doneRecordsCh, procResultCh, procFunc)

	//log proc result
	doneRecordsLogWait := make(chan bool)
	procResultLogWait := make(chan bool)
	go utils.Log(doneRecordsCh, cli.DoneRecordsPath, doneRecordsLogWait)
	go utils.Log(procResultCh, cli.ProcResultsPath, procResultLogWait)
	log.Println("Done Record in file: ", cli.DoneRecordsPath)
	log.Println("Proc Result in file: ", cli.ProcResultsPath)
	<-doneRecordsLogWait
	<-procResultLogWait

	log.Println("List and Proc done!")
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
