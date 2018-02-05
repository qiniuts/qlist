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

type batchProcFunc func(failedCh chan string, keys []string, cfg config.Config)

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

func (c *Client) Proc(recordsCh, processedCh, retCh chan string, batch batchProcFunc) {
	wg := sync.WaitGroup{}
	wg.Add(c.WorkerCount)

	for i := 0; i < c.WorkerCount; i++ {
		go c.worker(recordsCh, processedCh, retCh, batch, &wg)
	}

	wg.Wait()
	close(processedCh)
	close(retCh)

	return
}

func (c *Client) worker(recordsCh, processedCh, retCh chan string, batch batchProcFunc, wg *sync.WaitGroup) {

	defer wg.Done()
	records := []string{}
	for {
		record, ok := <-recordsCh
		if ok {
			processedCh <- record
			records = append(records, record)
			if len(records) < 1000 {
				continue
			}
		} else if len(records) == 0 {
			break
		}

		batch(retCh, records, c.Config)
		records = []string{}
	}
}
