package primitive

import (
	"fmt"
	"testing"
)

// Reverse
func TestReverse(t *testing.T) {
	list := new(List)
	strarr := []string{"test", "another", "reverse"}
	for i := range strarr {
		list.data = append(list.data, strarr[i])
	}
	ret := list.Reverse()
	fmt.Printf("TestReverse: %v\n", ret)
	if ret.data[0] != "reverse" {
		t.Error("Error Reverse()")
	}
}

// Map
// Basically return the underlying data
func TestMap(t *testing.T) {
	list := new(List)
	elem1 := interface{}("test")
	elem2 := interface{}(1)
	elem3 := interface{}(2.4)
	elem4 := interface{}(&elem3)
	list.data = append(list.data, elem1)
	list.data = append(list.data, elem2)
	list.data = append(list.data, elem3)
	list.data = append(list.data, elem4)
	ret := list.Map()
	fmt.Printf("%v\n", ret)
	if len(ret) != len(list.data) {
		t.Errorf("Error Map()")
	}
}

// Filter
// filter function => pair numbers
func TestFilter(t *testing.T) {
	list := new(List)
	intarr := []int{1, 3, -4, 6, 8, 7, -3}
	for i := range intarr {
		list.data = append(list.data, intarr[i])
	}
	ret := list.Filter(func(in interface{}) interface{} {
		val := in.(int)
		if val%2 == 0 {
			retInt := interface{}(val)
			return retInt
		}
		return nil
	})
	fmt.Printf("%v\n", ret)
	if len(ret.data) != 3 {
		t.Errorf("Error Filter() %v\n", ret)
	}
}

// Foldleft
// foldleft function => strings
func TestFoldLeft(t *testing.T) {
	list := new(List)
	strarr := []string{"fold", "this", "strings"}
	for i := range strarr {
		list.data = append(list.data, strarr[i])
	}
	ret := list.FoldLeft(func(in []interface{}) interface{} {
		var out string
		for i := range in {
			str := in[i].(string)
			out += str
		}
		ret := interface{}(out)
		return ret
	})
	fmt.Printf("%v\n", ret)
	if ret.data[0] != "foldthisstrings" {
		t.Errorf("Error FoldLeft() %v\n", ret)
	}
}
