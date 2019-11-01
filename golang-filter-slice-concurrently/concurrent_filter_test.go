package filter_test

import (
	"github.com/stretchr/testify/assert"
	filter "gitub.com/inge4pres/blog/concurrent-filter"
	"math/rand"
	"testing"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var oneKStrings = randStringSlice(1000)
var tenKStrings = randStringSlice(10000)
var hundredKStrings = randStringSlice(100000)

func randStringSlice(length int) []string {
	ret := make([]string, 0)
	for l := 0; l < length; l++ {
		vl := r.Intn(10) + 1
		b := make([]byte, vl)
		for i := range b {
			b[i] = letters[r.Intn(len(letters))]
		}
		ret = append(ret, string(b))
	}
	return ret
}

func Test_FilteredSliceIsLessThanOriginal(t *testing.T) {
	test := randStringSlice(4)
	toDel := "aaaaaaaaa"
	test = append(test, toDel)

	filtered := filter.Concurrently(test, filter.LenBiggerThan4)
	assert.True(t, len(filtered) <= 4)
	assert.NotContains(t, filtered, toDel)
}

func Test_FilteredUnoptmizedSliceIsLessThanOriginal(t *testing.T) {
	test := randStringSlice(4)
	toDel := "aaaaaaaaa"
	test = append(test, toDel)

	filtered := filter.ConcurrentlyUnoptimized(test, filter.LenBiggerThan4)
	assert.True(t, len(filtered) <= 4)
	assert.NotContains(t, filtered, toDel)
}

func Test_FilteredUnoptmizedSliceDoNotExcludeValid(t *testing.T) {
	test := []string{"a", "B", "c", "D"}
	toDel := "aaaaaaaaa"
	test = append(test, toDel)

	filtered := filter.ConcurrentlyUnoptimized(test, filter.LenBiggerThan4)
	assert.True(t, len(filtered) == 4)
	assert.NotContains(t, filtered, toDel)
}

func Benchmark_Concurrent_Unoptmized_1k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.ConcurrentlyUnoptimized(oneKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Concurrent_1k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Concurrently(oneKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Sequential_1k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Sequentially(oneKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Concurrent_Unoptmized_10k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.ConcurrentlyUnoptimized(tenKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Concurrent_10k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Concurrently(tenKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Sequential_10k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Sequentially(tenKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Concurrent_Unoptmized_100k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.ConcurrentlyUnoptimized(hundredKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Concurrent_100k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Concurrently(hundredKStrings, filter.LenBiggerThan4)
	}
}

func Benchmark_Sequential_100k(b *testing.B) {
	b.N = 1000
	for i := 0; i < b.N; i++ {
		filter.Sequentially(hundredKStrings, filter.LenBiggerThan4)
	}
}
