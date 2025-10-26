package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 { // слушать нечего
		return nil
	}
	if len(channels) == 1 { // если 1 то его и возвращаем
		return channels[0]
	}
	resChan := make(chan interface{})
	go func() {
		defer close(resChan)
		switch len(channels) {
		case 2: // если 2 то слушаем 2 канала только
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: // если больше 2-х то рекурсивно склеиваем половины
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()
	return resChan
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}