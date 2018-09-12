package utils

import (
	"math/rand"
	"time"
)

func Retry(f func() error) error {
	return retry(3, time.Second*20, f)
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if attempts <= 0 {
			return err
		}
		attempts--

		jitter := time.Duration(rand.Int63n(int64(sleep)))
		sleep = sleep + jitter/2
		time.Sleep(sleep)

		return retry(attempts, sleep, f)
	}

	return nil
}
