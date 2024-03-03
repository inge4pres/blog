package embedding

import (
	"fmt"
	"time"
)

// BusinessTask defines the interface for an operation that we want to instrument.
type BusinessTask interface {
	Execute() string
}

// BusinessLogic represents a product task to be run.
type BusinessLogic struct {
	Name string
}

// Execute implements the Execute method of the BusinessTask interface.
func (o BusinessLogic) Execute() string {
	return fmt.Sprintf("Executing %s operation", o.Name)
}

// InstrumentedTask is a decorator that adds observability to an operation.
type InstrumentedTask struct {
	BusinessTask

	latencyUsec,
	operationsCount int64
}

// GetMetrics returns the number of operations and the cumulative latency in microseconds.
func (io *InstrumentedTask) GetMetrics() (int64, int64) {
	return io.operationsCount, io.latencyUsec
}

// Execute implements the BusinessTask interface
// and adds observability by measuring how many times we ran it
// and keep track of latency through a counter.
func (io *InstrumentedTask) Execute() string {
	start := time.Now()
	defer func() {
		io.operationsCount += 1
		io.latencyUsec += time.Since(start).Microseconds()
	}()

	return io.BusinessTask.Execute()
}
