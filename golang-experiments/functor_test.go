package experiments

import "testing"

func TestFunctor(t *testing.T) {
	args := []interface{}{"Hello World", 2, true}
	returnedItem := SliceFunctor(lastItem, args)
	switch returnedItem.(type) {
	case bool:
		return
	default:
		t.Errorf("expected bool, found %T", returnedItem)
	}
}
