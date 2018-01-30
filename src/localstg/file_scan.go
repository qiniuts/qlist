package localstg

import (
	"os"
	"bufio"
	"log"
)

func List(inCh chan string, fpath string)  {

	fh, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}

	sc := bufio.NewScanner(fh)
	for sc.Scan() {
		inCh <- sc.Text()
	}

	if err := sc.Err(); err != nil {
		log.Println(err)
	}
	close(inCh)
}

