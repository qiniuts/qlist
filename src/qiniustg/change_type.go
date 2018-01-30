package qiniustg

import (
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"config"
)

//func ChangeFileType(ak, sk, bucket string, keyCh, doneCh, failedCh chan string, wg *sync.WaitGroup)  {
//	defer wg.Done()
//
//	mac := qbox.NewMac(ak, sk)
//	cfg := storage.Config{
//		UseHTTPS: false,
//	}
//	bucketManager := storage.NewBucketManager(mac, &cfg)
//
//	chtypeKeys := []string{}
//	chtypeOps := []string{}
//	fileType := 1
//
//	for {
//		key, ok := <- keyCh
//		if ok {
//
//			//receive and build your params
//			doneCh <- key
//			chtypeKeys = append(chtypeKeys, key)
//			chtypeOps = append(chtypeOps, storage.URIChangeType(bucket, key, fileType))
//			if  len(chtypeOps) < 1000 {
//				continue
//			}
//		} else if len(chtypeOps) == 0 {
//			break
//		}
//
//		//proc you params
//		rets, err := bucketManager.Batch(chtypeOps)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		for i, ret := range rets {
//			errStr := fmt.Sprintf("%d\t%s\t%s", ret.Code, chtypeKeys[i],ret.Data.Error)
//			failedCh <- errStr
//		}
//
//		chtypeOps = []string{}
//		chtypeKeys = []string{}
//	}
//}

func ChangeFileType2(failedCh chan string, keys []string, cfg config.Config) error  {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	stgCfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &stgCfg)
	chtypeOps := []string{}

	for _, key := range keys {
		chtypeOps = append(chtypeOps, storage.URIChangeType(cfg.Bucket, key, 1))
	}

	rets, err := bucketManager.Batch(chtypeOps)
	for i, ret := range rets {
		errStr := fmt.Sprintf("%d\t%s\t%s", ret.Code, keys[i],ret.Data.Error)
		failedCh <- errStr
	}

	return err
}
