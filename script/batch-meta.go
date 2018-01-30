package script

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"os"
	"strings"
	"sync"
)

func main() {
	var file string
	var bucket string
	var worker int
	var accessKey string
	var secretKey string
	flag.StringVar(&file, "file", "", "file list")
	flag.StringVar(&bucket, "bucket", "", "bucket name")
	flag.StringVar(&accessKey, "ak", "", "access key")
	flag.StringVar(&secretKey, "sk", "", "secret key")
	flag.IntVar(&worker, "worker", 1, "go routine count")
	flag.Parse()

	if accessKey == "" || secretKey == "" || bucket == "" || file == "" {
		fmt.Println("Usage: use -h to see usage")
		return
	}

	BatchChgm(accessKey, secretKey, bucket, file, worker)
}

const (
	BATCH_ALLOW_MAX = 10
)

func doBatchOperation(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func BatchChgm(accessKey, secretKey, bucket, file string, worker int) {
	mac := qbox.NewMac(accessKey, secretKey)

	var batchTasks chan func()
	var initBatchOnce sync.Once

	batchWaitGroup := sync.WaitGroup{}
	initBatchOnce.Do(func() {
		batchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doBatchOperation(batchTasks)
		}
	})

	fp, err := os.Open(file)
	if err != nil {
		fmt.Println("Open key mime map file error")
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)

	bucketManager := storage.NewBucketManager(mac, nil)
	chgmOps := make([]string, 0, BATCH_ALLOW_MAX)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "\t")
		if len(items) == 2 {
			key := items[0]
			chgmOp := URIChangeMeta2(bucket, key)
			chgmOps = append(chgmOps, chgmOp)
		} else {
			fmt.Println("Err:", line)
		}
		if len(chgmOps) == BATCH_ALLOW_MAX {
			toChgmOps := make([]string, len(chgmOps))
			copy(toChgmOps, chgmOps)

			batchWaitGroup.Add(1)
			batchTasks <- func() {
				defer batchWaitGroup.Done()
				batchChgm(bucketManager, toChgmOps)
			}
			chgmOps = make([]string, 0, BATCH_ALLOW_MAX)
		}
	}
	if len(chgmOps) > 0 {
		toChgmOps := make([]string, len(chgmOps))
		copy(toChgmOps, chgmOps)

		batchWaitGroup.Add(1)
		batchTasks <- func() {
			defer batchWaitGroup.Done()
			batchChgm(bucketManager, toChgmOps)
		}
	}

	batchWaitGroup.Wait()

}

func batchChgm(manager *storage.BucketManager, chgmOps []string) {
	rets, err := manager.Batch(chgmOps)
	if err != nil {
		// 遇到错误
		if _, ok := err.(*rpc.ErrorInfo); ok {
			for _, ret := range rets {
				// 200 为成功
				fmt.Printf("%d", ret.Code)
				if ret.Code != 200 {
					fmt.Printf("\t%s", ret.Data.Error)
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("batch error, %s\n", err)
		}
	} else {
		// 完全成功
		for _, ret := range rets {
			// 200 为成功
			fmt.Printf("%d", ret.Code)
			if ret.Code != 200 {
				fmt.Printf("\t%s", ret.Data.Error)
			}
			fmt.Println()
		}
	}

}

func URIChangeMeta2(bucket, key string) string {

		uri :=  fmt.Sprintf("/chgm/%s/x-qn-meta-!Content-Type/%s",
		base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", bucket, key))),
		base64.URLEncoding.EncodeToString([]byte("application/octet-stream")))

	return uri
}

//chgm/<EncodedEntryURI>/X-Qn-Meta-!Content-Type/<encode(application/octet-stream)
func URIChangeMeta(bucket, key string, etag, lastModified, md5 string) string {
	if md5 == "md5null" {
		return fmt.Sprintf("/chgm/%s/x-qn-meta-!ETag/%s/x-qn-meta-!Last-Modified/%s",
			base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", bucket, key))),
			base64.URLEncoding.EncodeToString([]byte(etag)), base64.URLEncoding.EncodeToString([]byte(lastModified)))
	} else {
		return fmt.Sprintf("/chgm/%s/x-qn-meta-!ETag/%s/x-qn-meta-!Last-Modified/%s/x-qn-meta-!X-CDO-Content-MD5/%s",
			base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", bucket, key))),
			base64.URLEncoding.EncodeToString([]byte(etag)), base64.URLEncoding.EncodeToString([]byte(lastModified)),
			base64.URLEncoding.EncodeToString([]byte(md5)))
	}
}
