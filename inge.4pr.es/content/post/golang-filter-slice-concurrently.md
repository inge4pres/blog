---
date: "2019-08-09T15:20:18+02:00"
title: "Golang-Filter-Slice-Concurrently"
author: "inge4pres"
emoji: true
comment: true

categories:
  - tech
  - work

tags:
  - golang
  - concurrency
  - programming

draft: true
---

Recently I was faced with the problem of filtering a collection in Go: given a slice of custom structs, how do you discard items matching a certain condition?
More over: can it be done in a concurrent-safe way? And if yes, is it faster than sequential filtering?

### Show me the code

TL;DR I found a way but I am not sure it's idiomatic: the core concept is to use a buffered channel to collect all items to be filtered, fan-out to multiple goroutines that will filter the items and fan back into a channel to concatenate back the filtered items.
Here's an example of code using `[]string` that can be adapted to custom structs.
 
```go
package filter 

type StringCmp func(string) bool

func Concurrently(str []string, filter StringCmp) []string {
	all := make(chan string, len(str))
	// send all values
	for _, s := range str {
		all <- s
	}
	// magic close, thanks Bill Kennedy!
	close(all)

	// build a partial results chainer and control channel to stop running goroutines
	// we also need a wg to let select know when we're done
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
```

### Is it good?

No! While at first glance the idea of executing in parallel the filtering of a slice, the cost of communication between memory and processor makes this choice a no-go for performance! 
See the benchmark yourself. 

```bash
go test -v -bench=. -benchmem ./...

goos: darwin
goarch: amd64
pkg: gitub.com/inge4pres/blog/concurrent-filter
Benchmark_Concurrent_Unoptmized_1k-8     	    1000	    550761 ns/op	   16974 B/op	      17 allocs/op
Benchmark_Concurrent_1k-8                	    1000	    213386 ns/op	   53119 B/op	      60 allocs/op
Benchmark_Sequential_1k-8                	    1000	     13567 ns/op	   16368 B/op	      10 allocs/op
Benchmark_Concurrent_Unoptmized_10k-8    	    1000	   5464960 ns/op	  253133 B/op	      23 allocs/op
Benchmark_Concurrent_10k-8               	    1000	   1905460 ns/op	  556386 B/op	      94 allocs/op
Benchmark_Sequential_10k-8               	    1000	    258709 ns/op	  252528 B/op	      16 allocs/op
Benchmark_Concurrent_Unoptmized_100k-8   	    1000	  58529722 ns/op	 3652848 B/op	      33 allocs/op
Benchmark_Concurrent_100k-8              	    1000	  19772039 ns/op	 6828590 B/op	     149 allocs/op
Benchmark_Sequential_100k-8              	    1000	   3446710 ns/op	 3652208 B/op	      26 allocs/op
PASS
ok  	gitub.com/inge4pres/blog/concurrent-filter	114.297s
```
