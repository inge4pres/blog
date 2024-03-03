package embedding

import (
	"fmt"
	"testing"
	"time"
)

func TestInstrumentedOperation(t *testing.T) {
	const (
		op    = "getProductX"
		count = 100
	)
	instrumented := &InstrumentedTask{BusinessTask: BusinessLogic{op}}
	for range count {
		instrumented.Execute()
		time.Sleep(1 * time.Microsecond)
	}
	execCount, latency := instrumented.GetMetrics()
	if execCount != count {
		t.Errorf("Expected %d, got %d", count, execCount)
	}
	if latency == 0 {
		t.Errorf("Expected a non-zero value for latency")
	}
}

func ExampleInstrumentedTask_Execute() {
	const op = "getProductX"
	instrumented := &InstrumentedTask{BusinessTask: BusinessLogic{op}}
	result := instrumented.Execute()
	// Output: Executing getProductX operation
	fmt.Println(result)
}
