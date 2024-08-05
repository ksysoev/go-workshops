package concurrency

import (
	"testing"
	"time"
)

// Go follows a model of concurrency called the fork-join model.
// This means that our program can split into it’s own execution branch to be run concurrently with its main branch.
// At some point in the future, the two branches of execution will be joined together again
//
// Main Task
// | (Fork)
// v
// +-----------+
// |           |
// Main Task   Child Task
// |           |
// +-----------+
// | (Join)
// v
// Final Result

func TestForkJoin(t *testing.T) {
	done := make(chan struct{})

	go func() {
		t.Log("Child goroutine")
	}()

	t.Log("Parent goroutine")

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Errorf("Expected call to return")
	}
}
