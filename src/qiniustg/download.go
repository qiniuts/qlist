package qiniustg

import (
	"config"
	"fmt"
	"github.com/cavaliercoder/grab"
	"utils"
)

func Download(recordsCh, retCh chan string, cfg config.Config) {

	for url1 := range recordsCh {

		fmt.Println(url1)

		err := utils.Retry(func() error {
			_, err := grab.Get(".", url1)
			return err
		})

		if err != nil {
			retCh <- fmt.Sprintf("%s\t%#v", url1, err)
		} else {
			retCh <- fmt.Sprintf("%s\t%#v", url1, err)
		}
	}
}
