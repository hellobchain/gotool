package gretry

import (
	"time"
)

func Do(max int, fn func() error) error {
	var err error
	for i := 0; i < max; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < max-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	return err
}
