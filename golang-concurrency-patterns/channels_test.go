package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestChannelIterateOverRange(t *testing.T) {
	numbers := make(chan int)
	go func(nums <-chan int) {
		for n := range nums {
			fmt.Println(printNumber(n))
		}
	}(numbers)
	for i := 0; i < nums; i++ {
		numbers <- i
	}
	close(numbers)
}

func TestChannelBuffered(t *testing.T) {
	numbers := make(chan int, 5)
	go func(nums <-chan int) {
		for n := range nums {
			fmt.Println(printNumber(n))
		}
	}(numbers)
	for i := 0; i < nums; i++ {
		go func(c int) {
			numbers <- c
		}(i)
	}
}

func TestMultipleChannelsSelect(t *testing.T) {
	numbers := make(chan int)
	greets := make(chan string)
	done := make(chan bool)
	finish := make(chan bool)
	// reading all in the background, exiting when done
	go func() {
		for {
			select {
			case n := <-numbers:
				fmt.Println(printNumber(n))
			case p := <-greets:
				fmt.Println(sayHello(p))
			case <-finish:
				return
			}
		}
	}()
	// sending numbers
	go func() {
		for i := 0; i < nums; i++ {
			numbers <- i
		}
		done <- true
	}()
	// sending strings
	go func() {
		for _, person := range people {
			greets <- person
		}
		done <- true
	}()

	// read to block until goroutines complete
	<-done
	<-done
	// exit from first goroutine
	finish <- true
}

func TestChannelRangeOnGoroutines(t *testing.T) {
	var wg sync.WaitGroup
	numbers := make(chan int)
	go func(nums <-chan int) {
		for n := range nums {
			fmt.Println(printNumber(n))
		}
	}(numbers)
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			numbers <- c
		}(i)
	}
	wg.Wait()
	close(numbers)
}
