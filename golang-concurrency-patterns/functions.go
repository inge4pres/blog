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

func numberConvert(num string) (int, error) {
	return strconv.Atoi(num)
}

func sayHello(who string) string {
	return fmt.Sprintf("Hello %s!", who)
}
