package experiments

import (
	"fmt"
	"testing"
)

var stringTypeFilterFn = func(input interface{}) bool {
	switch t := input.(type) {
	case string:
		return true
	default:
		fmt.Printf("detected type %T", t)
		return false
	}
}

func TestStream(t *testing.T) {
	underTest := new(Iterable)
	underTest.list = []interface{}{"Hi!", -4, false}
	result := underTest.Stream().Collect(stringTypeFilterFn).Close()
	if len(result) != 1 {
		t.Errorf("expected to collect just 1 item, found %v", result)
	}
}
