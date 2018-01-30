package script

import (
	//"qiniu.com/auth/qiniumac"
	"qiniu.com/auth/qiniumac.v1"
	"net/http"
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"encoding/json"
	"time"
	"os"
	"bufio"
	"strings"
	"flag"
	"strconv"
)

func main() {

	ak := flag.String("ak", "your_access_key", "")
	sk := flag.String("sk", "your_secret_key", "")
	start_ := flag.String("start", "start_date", "")
	end_ := flag.String("end", "end_date", "")
	streamFile := flag.String("stream", "stream_file_path", "")
	flag.Parse()


	//start := "1513440000"
	//end := "1514736000"

	startDate, err := time.Parse("2006-01-02", *start_)
	if err != nil {
		panic(err)
	}
	start := strconv.Itoa(int(startDate.Unix()))

	endDate, err := time.Parse("2006-01-02", *end_)
	if err != nil {
		panic(err)
	}
	end := strconv.Itoa(int(endDate.Unix()))


	streams, err := readLines(*streamFile)
	if err != nil {
		panic(err)
	}
	//streams := []string{"bb51ef892bfdefaa29351276bfa47ba9"}


	for _, stream := range streams {

		encodedStreamKey := base64.URLEncoding.EncodeToString([]byte(stream))

		u := "http://pili.qiniuapi.com/v2/hubs/live_panda/streams/" + encodedStreamKey + "/historyinfo?start=" + start + "&end=" + end
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			panic(err)
		}
		su := "1380637639/0"

		sign, err := qiniumac.DefaultRequestSigner.SignAdmin([]byte(*sk), req, su)
		if err != nil {
			panic(err)
		}

		//<SuInfo>:<AK>:<Sign>
		auth := "QiniuAdmin " + su + ":" + *ak + ":" + base64.URLEncoding.EncodeToString(sign)

		req.Header.Set("Authorization", auth)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		ret := struct {
			Items []Stream `json:"items"`
		}{}
		json.Unmarshal(body, &ret)


		for _, item := range ret.Items {

			t := time.Unix(item.Time, 0).Format("2006-01-02T15:04:05")

			fmt.Printf("%s\t%d\t%s\t%d\n", stream, item.Time, t,item.Bps)
		}
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

type Stream struct {
	Time int64 `json:"time"`
	Bps int `json:"bps"`
}
