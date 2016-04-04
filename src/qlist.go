package main

import(
	"qiniupkg.com/api.v7/kodo"

	"time"
	"io"
	"strings"
	"encoding/base64"
	"net/url"
	"sync"
	"os"
	"bufio"
	"fmt"
	"strconv"
	"net/http"
	"crypto/hmac"
	"crypto/sha1"
	"flag"
)



func main() {

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: qlist <prefix> <sleep_ms> [<marker>]")
		return
	}

	cfg := &kodo.Config{}
	cfg.AccessKey = ""
	cfg.SecretKey = ""

	cli := kodo.New(0, cfg)


	bucketName := "babytree"
	bucket := cli.Bucket(bucketName)

	prefix := flag.Arg(0)
	ms, _ := strconv.Atoi(flag.Arg(1))
	deli := ""
	marker := flag.Arg(2)
	limit := 300

	fname := bucketName + "_" + base64.URLEncoding.EncodeToString([]byte(prefix))
	defer func() {
		fmt.Println("marker ==> ", marker)
	}()


	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fp := wd + "/" + fname
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}
	w := bufio.NewWriter(f)
	defer func () {
		w.Flush()
		f.Sync()	
		f.Close()
	}()

	logp := wd + "/" + fname + "_err.log"
	lf, err := os.OpenFile(logp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
	    panic(err)
	}
	lw := bufio.NewWriter(lf)
	defer func () {
		lw.Flush()
		lf.Sync()	
		lf.Close()
	}()



    resps := make(chan string)
   	items := make(chan []kodo.ListItem)	

    var wg sync.WaitGroup
    domain := "http://7u2pvv.com0.z0.glb.qiniucdn.com/"

    go func() {
        for response := range resps {
            if strings.HasPrefix(response, "200 OK") {
            	fmt.Fprintln(w, response[7:])
	            fmt.Println(response[7:])
            } else {
            	fmt.Fprintln(lw, response)
            }
        }
    }()

	for {
		entries, _, marker_, err := bucket.List(nil, prefix, deli, marker, limit);
		if err != nil && err != io.EOF {
			fmt.Println("Error", err)	
			return 
		}
		marker = marker_
		time.Sleep(time.Duration(ms) * time.Millisecond)



		items <- entries

		if err == io.EOF { 
			fmt.Println("done <===========>")	
			break
		}
	}
	wg.Wait()
}

func Proc(items chan []kodo.ListItem, concurrency int) {



    wg.Add(len(entries))
	for i := 0; i < concurrency; i++ {
		go func(kodo.ListItem) {
			defer wg.Done()	

			if !isImage(e) {
				return
			}

			url := domain + e.Key
			url = saveasUrl(url + "?imageView2/2/w/2560/h/2560/q/95/format/webp", cfg.AccessKey, []byte(cfg.SecretKey), bucketName, e.Key + "_webp")
			fmt.Println(url)
            res, err := http.Head(url)
            if err != nil {
            	fmt.Fprintln(lw, err)
            } else {
                resps <- string(res.Status + "\t" + e.Key)
            }
		}(e)
	}
}


func isImage(e kodo.ListItem) bool {
	fmt.Println("--->", e.Key)
	return strings.HasPrefix(e.MimeType, "image") && !strings.HasSuffix(e.Key, "_webp")
}

func saveasUrl(url_, accessKey string, secretKey []byte, bucket, key string) string {

      encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(bucket+":"+key))

      url_ += "|saveas/" + encodedEntryURI

      h := hmac.New(sha1.New, secretKey)

      u, _ := url.Parse(url_)
      io.WriteString(h, u.Host + u.RequestURI())

      d := h.Sum(nil)
      sign := accessKey + ":" + base64.URLEncoding.EncodeToString(d)

      return url_ + "/sign/" + sign

}



