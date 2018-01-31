package utils

import (
	"os"
	"bufio"
)

func FileLastLine(fpath string) (ll []byte, err error) {
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

func Log(logCh chan string, fpath string, done chan bool)  {
	fh, err := os.OpenFile(fpath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	for key := range logCh {
		fh.WriteString(key + "\n")
	}
	close(done)
}

