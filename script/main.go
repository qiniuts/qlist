package script

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	xlog "github.com/qiniu/xlog.v1"
	"github.com/sbunce/bson"
	gbson "labix.org/v2/mgo/bson"
	ebdtypes "qbox.us/ebd/api/types"
	"qbox.us/fh/fhver"

	_ "qbox.us/autoua"
	_ "qbox.us/profile"
)

type Fids []uint64

var fidmap = make(map[uint64]bool, 200000000)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
	fidFiles := flag.String("fids", "fids", "fids")
	process := flag.Int("process", 8, "process")
	out := flag.String("out", "out", "out")
	flag.Parse()
	f, err := os.Open(*fidFiles)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	var fid uint64
	var fidCount int
	for {
		if fidCount%100000 == 0 {
			log.Println("load fid", fidCount)
		}
		fidCount++
		_, err := fmt.Fscanf(f, "%d", &fid)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Panic(err)
		}
		fidmap[fid] = true
	}
	log.Println("fid count", fidCount)
	filterCh := make(chan []byte, 100)
	writeCh := make(chan []byte, 100)
	doneCh := make(chan bool)
	go readFromStdin(filterCh)
	wg := &sync.WaitGroup{}
	wg.Add(*process)
	for i := 0; i < *process; i++ {
		go filter(filterCh, writeCh, wg)
	}
	go writerToFile(writeCh, *out, doneCh)
	wg.Wait()
	close(writeCh)
	<-doneCh
}

func readFromStdin(ch chan<- []byte) {
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := bson.ReadOne(reader)
		if err != nil {
			if err == io.EOF {
				close(ch)
				break
			}
			log.Panic(err)
		}
		ch <- b
	}
}

func filter(in <-chan []byte, out chan<- []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	var doc Doc
	xl := xlog.NewDummy()
	for b := range in {
		err := gbson.Unmarshal(b, &doc)
		if err != nil {
			log.Panic(err)
		}
		ver := fhver.FhVer(doc.Fh)
		if ver < 5 || ver > 7 {
			continue
		}
		fhi, err := ebdtypes.DecodeFh(doc.Fh)
		if err != nil {
			m, _ := bson.BSON(b).Map()
			log.Panicf("%v %+v", err, m)
		}
		if _, ok := fidmap[fhi.Fid]; ok {
			xl.Infof("%v\t%v\t%#v", fhi.Gid, fhi.Fsize, doc.ID)
			out <- b
		}
	}
}

func writerToFile(in <-chan []byte, filename string, done chan<- bool) {
	f, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	var count int64
	for b := range in {
		if count%10000 == 0 {
			log.Println("count", count)
		}
		_, err = f.Write(b)
		if err != nil {
			log.Panic(err)
		}
		count += 1
	}
	err = f.Close()
	if err != nil {
		log.Panic(err)
	}
	close(done)
	log.Println("total", count)
}

type Doc struct {
	ID string `bson:"_id"`
	Fh []byte `bson:"fh"`
}