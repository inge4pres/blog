package filter_test

import (
	"testing"

	filter "github.com/inge4pres/golang-functional-patterns"
	"github.com/stretchr/testify/assert"
)

var underTest = []string{"a", "b", "c", "d", "e", "aaa", "f"}

func Test_filterObject(t *testing.T) {
	o := filter.Obj{
		Strings: underTest,
	}
	o.CriteriaFactory().WithLengthBiggerThan(2)
	o.Filter()

	assert.Equal(t, 6, len(o.Strings))
	assert.Equal(t, "f", o.Strings[5])
}

func Test_filterFunctional(t *testing.T) {
	f := filter.Strings(underTest, filter.BiggerThan(2))

	assert.Equal(t, 7, len(f))
	assert.Equal(t, "f", f[5])
}

func Benchmark_filterObject(b *testing.B){
	b.N = 10000000
	for i:=0; i<b.N; i++ {
		o := filter.Obj{
			Strings: underTest,
		}
		o.CriteriaFactory().WithLengthBiggerThan(2)
		o.Filter()
	}
}

func Benchmark_filterFunctional(b *testing.B){
	b.N = 10000000
	for i:=0; i<b.N; i++ {
		filter.Strings(underTest, filter.BiggerThan(2))
	}
}