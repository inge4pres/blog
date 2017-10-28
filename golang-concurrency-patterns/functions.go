package concurrency

import (
	"fmt"
	"math/rand"
	"strconv"
)

func randomInt() int {
	return rand.Int()
}

func printNumber(val int) string {
	return fmt.Sprintf("%s", strconv.Itoa(val))
}

func sayHello(who string) string {
	return fmt.Sprintf("Hello %s!", who)
}
