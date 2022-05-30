package func_test

import (
	"sync"
	"time"
)

/*
   EXAMPLE 1.

   This example is using only one goroutine. It creates CPU overhead, because there are many nested if statements and for loops, as well as select statement with continue statements in it. The approach in this example is shown for educational purposes only and should not be used in real life cases.
*/

func all(b []bool) bool {
	for _, v := range b {
		if !v {
			return false
		}
	}
	return true
}

func combine[T any](chans ...chan T) (out chan T) {
	out = make(chan T)
	go func() {
		closed := make([]bool, len(chans))
		for {
			for n, chann := range chans {
				if !closed[n] {
					select {
					case v, ok := <-chann:
						if !ok {
							closed[n] = true
							continue
						}
						out <- v
					default:
						continue
					}
				}
			}
			if all(closed) {
				close(out)
				return
			}
		}
	}()
	return
}

/*
 EXAMPLE 2.

 This function is much better than function in first example. It uses one goroutine for one channel, and one goroutine for the WaitGroup monitoring. So, this function will create `len(channels) + 1` goroutines. Here is no CPU overhead while using this function
*/

func combineTwo[T any](chans ...chan T) (out chan T) {
	out = make(chan T)
	wg := new(sync.WaitGroup)

	wg.Add(len(chans))
	for n, ch := range chans {
		go func(x int, c chan T) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(n, ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func sendData[T any](start, stop int,
	sleep time.Duration, sndVal T, outCh chan<- T) {

	for i := start; i < stop; i++ {
		time.Sleep(sleep)
		outCh <- sndVal
	}
	close(outCh)
}

func consume[T any](inChan <-chan T) (x []T) {
	x = make([]T, 0)
	for v := range inChan {
		x = append(x, v)
	}

	return x
}
