package concurrency

import (
	"fmt"
	"testing"
)

func TestUnbufferedChannel(t *testing.T) {
	numbers := make(chan int)
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

func TestBufferedChannel(t *testing.T) {
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

func TestChannelWaitForFinishOnMultipleChannels(t *testing.T) {
	numbers := make(chan int, 10)
	greets := make(chan string)
	done := make(chan bool)
	// reading all in the background, exiting when done
	go func() {
		for {
			select {
			case n := <-numbers:
				fmt.Println(printNumber(n))
			case p := <-greets:
				fmt.Println(sayHello(p))
			}
		}
	}()
	// sending numbers
	go func() {
		for i := 0; i < nums; i++ {
			go func(c int) {
				numbers <- c
			}(i)
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

	// wait reading from done
	<-done
	<-done
}
