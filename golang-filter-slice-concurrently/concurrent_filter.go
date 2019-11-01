package filter

import (
	"runtime"
	"sync"
)

type StringCmp func(string) bool

func Concurrently(str []string, filter StringCmp) []string {
	all := make(chan string, len(str))
	// send all values
	for _, s := range str {
		all <- s
	}
	// magic close, thanks Bill Kennedy!
	close(all)

	// build a partial result chainer and control chan to stop running goroutines
	// we also need a wg to let the select know when we're done
	done := make(chan bool)
	workers := runtime.NumCPU()
	valid := make(chan []string, workers)
	processing := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		processing.Add(1)
		go filterGroup(all, valid, filter)
	}
	// create the output chained result slice
	filtered := make([]string, 0)
	// listen for the fan-in items or finish when done
	go func() {
		for {
			select {
			case <-done:
				return
				// chain when getting data from chan
			case partial := <-valid:
				filtered = append(filtered, partial...)
				// and signal
				processing.Done()
			}
		}
	}()
	// ready to close
	processing.Wait()
	done <- true
	close(valid)
	return filtered
}

func ConcurrentlyUnoptimized(str []string, filter StringCmp) []string {
	// build a result chainer and control chan to stop running goroutines
	// we also need a wg to let the select know when we're done
	done := make(chan bool)
	workers := runtime.NumCPU()
	all := make(chan string, workers)
	valid := make(chan string, workers)
	processing := &sync.WaitGroup{}

	// start background goroutines to read from the channel,
	// filter items and send back to output channel
	for i := 0; i < workers; i++ {
		processing.Add(1)
		go filterSingle(all, valid, filter, processing)
	}

	// create the output chained result slice
	filtered := make([]string, 0)
	// listen for the fan-in items or finish when done
	go func() {
		for {
			select {
			case <-done:
				return
			case s := <-valid:
				filtered = append(filtered, s)
			}
		}
	}()
	// stream values in!
	for _, s := range str {
		all <- s
	}
	// ready to close, shut down input so goroutines will signal completion
	close(all)
	processing.Wait()
	// signal the readers we're done
	done <- true
	close(valid)
	// signal the sender to close
	return filtered
}

func Sequentially(str []string, filter StringCmp) []string {
	ret := make([]string, 0)
	for _, s := range str {
		if filter(s) {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}

func filterGroup(in <-chan string, out chan<- []string, filterOutFn StringCmp) {
	partial := make([]string, 0)
	for s := range in {
		if filterOutFn(s) {
			continue
		}
		partial = append(partial, s)
	}
	out <- partial
}

func filterSingle(in <-chan string, out chan<- string, filterOutFn StringCmp, g *sync.WaitGroup) {
	for s := range in {
		if filterOutFn(s) {
			continue
		}
		out <- s
	}
	g.Done()
}

func LenBiggerThan4(in string) bool {
	return len(in) > 4
}
