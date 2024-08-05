package concurrency

import (
	"testing"
	"time"
)

//

func TestForkJoin(t *testing.T) {
	done := make(chan struct{})

	go func() {
		t.Log("Child goroutine")
		close(done)
	}()

	t.Log("Parent goroutine")

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Errorf("Expected call to return")
	}
}
