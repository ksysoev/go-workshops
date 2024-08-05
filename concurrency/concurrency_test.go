package concurrency

import (
	"context"
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
		t.Errorf("Expected to done channel to be closed")
	}
}

// The main should be able to control the child goroutine execution.
// For example if the main goroutine will be requested to stop, it should be able to cancel the child goroutine execution.
// This can be done by using context package.
func TestParrentControl(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		t.Log("Child goroutine")
		<-time.After(10 * time.Second)

		close(done)
	}()

	t.Log("Parent cancel child goroutine execution")
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Errorf("Expected to done channel to be closed")
	}
}
