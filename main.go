package main

import (
	"config"
	"flag"
	"fmt"
	"localstg"
	"log"
	"qiniustg"
	"runtime"
	"sync"
	"utils"
)

type procFunc func(recordsCh, retCh chan string, cfg config.Config)

var procFuncs = map[string]procFunc{
	"req":         qiniustg.HttpReq,
	"download":    qiniustg.Download,
	"bucket_list": localstg.BucketList,
	"adl":         qiniustg.Adl,
	"chstatus":    batchFunc(qiniustg.ChangeFileStatus),
	"chtype":      batchFunc(qiniustg.ChangeFileType),
	"async_fetch": qiniustg.AsyncFetch,
	"md5":         qiniustg.Md5,
}

func usage() {
	fmt.Println(
		`
		./qlist -cfg_path cfg.json req
		./qlist -cfg_path cfg.json download
		./qlist -cfg_path cfg.json bucket_list
		./qlist -cfg_path cfg.json chstatus
		./qlist -cfg_path cfg.json chtype
		./qlist -cfg_path cfg.json async_fetch
		./qlist -cfg_path cfg.json md5
	`)
}

func main() {

	cfgPath := flag.String("cfg_path", "", "")
	flag.Parse()

	if *cfgPath == "" {
		usage()
		return
	}

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

	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	//channels to cache records and proc result
	recordsCh := make(chan string, 1000*100)
	procResultCh := make(chan string, 1000*10)

	if cfg.IsQiniuSrc() {
		//list records qiniu storage
		go qiniustg.NewClient(cfg).List(recordsCh)
	} else {
		//list records local file
		go localstg.NewClient(cfg).List(recordsCh)
	}

	//proc records
	go cli.work(recordsCh, procResultCh, procFunc)

	//log proc result
	procResultLogWait := make(chan bool)
	go utils.Log(procResultCh, cli.ProcResultsPath, procResultLogWait)
	log.Println("Proc Result in file: ", cli.ProcResultsPath)
	<-procResultLogWait

	log.Println("List and Proc done!")
}

type Client struct {
	config.Config
}

func (c *Client) work(recordsCh, retCh chan string, proc procFunc) {
	wg := sync.WaitGroup{}
	if c.WorkerCount == 0 {
		c.WorkerCount = 10
	}
	wg.Add(c.WorkerCount)

	for i := 0; i < c.WorkerCount; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			proc(recordsCh, retCh, c.Config)
		}(&wg)
	}

	wg.Wait()
	close(retCh)
	return
}

func batchFunc(batch func(retCh chan string, records []string, cfg config.Config)) func(recordsCh, retCh chan string, cfg config.Config) {
	return func(recordsCh, retCh chan string, cfg config.Config) {
		records := []string{}
		for {
			record, ok := <-recordsCh
			if ok {
				records = append(records, record)
				if len(records) < 1000 {
					continue
				}
			}

			if len(records) == 0 {
				break
			}
			batch(retCh, records, cfg)
			records = []string{}
		}
	}
}
