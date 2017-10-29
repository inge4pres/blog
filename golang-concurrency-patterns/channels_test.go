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
	numbers := make(chan string, 5)
	go func(nums <-chan string) {
		for n := range nums {
			fmt.Println(n)
		}
	}(numbers)
	for i := 0; i < nums; i++ {
		go func(c int) {
			numbers <- printNumber(c)
		}(i)
	}
}

func TestMultipleChannelsSelect(t *testing.T) {
	numbers := make(chan string)
	greets := make(chan string)
	done := make(chan bool)
	finish := make(chan bool)
	// reading all in the background, exiting when done
	go func() {
		for {
			select {
			case n := <-numbers:
				fmt.Println(n)
			case p := <-greets:
				fmt.Println(p)
			case <-finish:
				return
			}
		}
	}()
	// sending numbers
	go func() {
		for i := 0; i < nums; i++ {
			numbers <- printNumber(i)
		}
		done <- true
	}()
	// sending strings
	go func() {
		for _, person := range people {
			greets <- sayHello(person)
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
	numbers := make(chan string)
	go func(nums <-chan string) {
		for n := range nums {
			fmt.Println(n)
		}
	}(numbers)
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()
			numbers <- printNumber(c)
		}(i)
	}
	wg.Wait()
	close(numbers)
}
