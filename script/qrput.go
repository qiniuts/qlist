package script

import (
	"fmt"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
	. "golang.org/x/net/context"
	"qiniupkg.com/api.v7/kodocli"
	cf "qiniupkg.com/api.v7/conf"
    "os"
	"path/filepath"
    "net/http"

	"path"
    "io/ioutil"
    "time"
)

func Help() {
	fmt.Print(`
Usage:
    qcheetah_up <fileName>
			`)
}

func main() {
	if len(os.Args) < 2 {
		Help()
		return
	}

    fmt.Println("Date Time:", time.Now().Format("2016-01-02 15:04:05 MST"))
    fmt.Println("TimeStamp:", time.Now().Unix())

    err := Put("na-test", os.Args[1])
	if err != nil {
		fmt.Println("\n", err)
	}
    ip()
}


func Put(scope, fileName string) (err error) {

	key := path.Base(fileName)

//	uptoken, err := getToken(scope)
//    fmt.Println(uptoken)
//	if err != nil {
//		return
//	}

    uptoken := "94uMw0jMtPC1JGoSiDJcWIbezbTnM8AchIaWDqfJ:HS4KX5xp2GFXlK3j2GXKHRL76AY=:eyJzY29wZSI6Im5hLXRlc3QiLCJkZWFkbGluZSI6MjEyMTMwMTcxNX0="

	uploader := kodocli.NewUploader(0, nil)
	extra := &kodocli.RputExtra{}

//	err = InitExtra(extra, fileName)
//	if err != nil {
//		return
//	}


    cf.SetAppName("94uMw0jMtPC1JGoSiDJcWIbezbTnM8AchIaWDqfJ")

    //设置分片大小
    v := &kodocli.Settings{
        Workers: 8,
        ChunkSize: 4*1024*1024,
    }
    kodocli.SetSettings(v)

	var ret interface{}
	ctx := TODO()

    begin := time.Now()
	err = uploader.RputFile(ctx, &ret, uptoken, key, fileName, extra)
	if err != nil {
		return
	}

    dur := time.Since(begin)
    sec := dur.Seconds()


    fi, err := os.Stat(fileName);
    if err != nil {
        fmt.Println("os.Stat", err)
    }
    size := fi.Size()


    fmt.Printf("File Size: %d Byte\n", size)
    fmt.Printf("Elapsed Time: %.2f Sec\n", sec)
    fmt.Printf("Upload Speed: %.2f KB/s\n", float64(size)/1024/sec)
    fmt.Println(ret)

	return
}

func InitExtra(extra *kodocli.RputExtra, fileName string) (err error) {


	fi, err := os.Lstat(fileName)
	if err != nil {
		return
	}

	// blockCnt
	blockCnt := kodocli.BlockCount(fi.Size())
	extra.Progresses = make([]kodocli.BlkputRet, blockCnt)


    blkNum := 0
	ShowProgress(blockCnt, blkNum, 50)

	var inc func()
	if blockCnt > 50 {
		onePointCnt := blockCnt / 50
		curCnt := blkNum % onePointCnt
		inc = func() {
			curCnt++
			if curCnt == onePointCnt {
				fmt.Print(">")
				curCnt = 0
			}
		}
	} else {
		oneCntPoint := 50 / blockCnt
		inc = func() {
			for i := 0; i < oneCntPoint; i++ {
				fmt.Print(">")
			}
		}
	}

	extra.Notify = func(blkIdx int, blkSize int, ret *kodocli.BlkputRet) {
		if uint32(blkSize) == ret.Offset {
			inc()
		}
		return
	}

	return
}

func ShowProgress(blockCnt, blkNum, pointNum int) {
	fmt.Print("|")
	for i := 0; i < pointNum; i++ {
		fmt.Print("=")
	}
	fmt.Print("|100%\r\n ")

	n := blkNum * pointNum / blockCnt
	for i := 0; i < n; i++ {
		fmt.Print(">")
	}
}

func getToken(scope string) (uptoken string, err error) {

	conf.ACCESS_KEY = "94uMw0jMtPC1JGoSiDJcWIbezbTnM8AchIaWDqfJ"

	putPolicy := rs.PutPolicy{
		Scope:   scope,
		Expires: 20 * 365 * 24 * 3600,
	}

	return putPolicy.Token(nil), nil
}


func GetLogDir() (dir string) {
	return filepath.Join(os.Getenv("HOME"), ".qrput")
}

func ip() {
    resp, _ := http.Get("http://myip.ipip.net")
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    fmt.Println(string(body))
}
