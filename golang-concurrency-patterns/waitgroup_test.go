package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestWaitGroupNoConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(ops)
	for i := 0; i < ops; i++ {
		printNumber(randomInt())
		wg.Done()
	}
	wg.Wait()
}

func TestWaitGroupFanoutFactor2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(ops)
	go func() {
		for i := 0; i < ops/2; i++ {
			defer wg.Done()
			printNumber(randomInt())
		}
	}()
	go func() {
		for i := ops / 2; i < ops; i++ {
			defer wg.Done()
			printNumber(randomInt())
		}
	}()
	wg.Wait()
}

func TestWaitGroupFanout(t *testing.T) {
	var wg sync.WaitGroup
	// can also be adding to the group the whole len(people) here like
	// wg.Add(len(people))
	for _, person := range people {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			fmt.Println(sayHello(p))
		}(person)
	}
	wg.Wait()
}
