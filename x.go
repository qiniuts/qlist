package main

import (
	"os"
	"fmt"
	"bufio"
)

func main()  {

	line, err := lastFileLine("log.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(".....>", string(line))
}

func lastFileLine(fpath string) (ll []byte, err error) {
	fi, err := os.Stat(fpath)
	if err != nil {
		return
	}
	fsize := fi.Size()

	fh, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer fh.Close()

	offset := fsize - 1
	buf := make([]byte, 1)
	for offset >= 0  {

		_, err = fh.ReadAt(buf, offset)
		if err != nil {
			return
		}

		if offset == 0 {
			ll, _, err = bufio.NewReader(fh).ReadLine()
			return
		}

		if string(buf) == "\n" && offset != fsize - 1{

			l := fsize - offset - 1
			lbuf := make([]byte, l)

			fh.ReadAt(lbuf, offset + 1)
			ll = lbuf
			return
		}
		offset--
	}

	return
}
