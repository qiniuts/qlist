package localstg

import (
	"os"
	"bufio"
)

func List(inCh chan string, fpath string)  {

	f, err := os.OpenFile(fpath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		inCh <- sc.Text()
	}

	close(inCh)
}

