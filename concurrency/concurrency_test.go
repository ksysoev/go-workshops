package concurrency

import (
	"context"
	"sync"
	"sync/atomic"
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
		t.Log("Child goroutine done")
	case <-time.After(time.Millisecond):
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
		t.Log("Child goroutine done")
	case <-time.After(time.Millisecond):
		t.Errorf("Expected to done channel to be closed")
	}
}

// Unbounded concurrency can lead to resource exhaustion and poor performance due to contention.
// To limit the number of goroutines that can run concurrently, we can use a semaphore.
// A semaphore is a synchronization primitive that limits the number of concurrent operations.
// It is used to control access to a shared resource.
// We can use a buffered channel to implement a semaphore.
func TestSemaphoreWithChannels(t *testing.T) {
	c := atomic.Int32{}
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			val := c.Add(1)
			defer c.Add(-1)

			time.Sleep(1 * time.Millisecond)

			if val > 3 {
				t.Error("Expected to have only 3 goroutines running concurrently")
			}
		}()
	}

	wg.Wait()
}
