package localstg

import (
	"bufio"
	"os"
)

func (c QNClient) List(inCh chan string) {

	f, err := os.OpenFile(c.ToDoRecordsPath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		inCh <- sc.Text()
	}

	close(inCh)
}
