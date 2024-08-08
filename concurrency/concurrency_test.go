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

// Data races and Race conditions, what is the difference?
// Data race is a condition in which two goroutines access the same variable concurrently and at least one of the accesses is a write.
// On the other hand Race condition is a condition in which the program’s output is dependent on the sequence or timing of uncontrollable events,
// which can lead to non-deterministic behavior.
func TestRaceCondition(t *testing.T) {
	data := 0

	go func() {
		data++
	}()

	if data != 1 {
		t.Error("Expected data to incremented by 1")
	}
}

// When something is considered atomic, or to have the property of atomicity,
// this means that within the context that it is operating, it is indivisible, or uninterruptible.
// Most statements in Go are not atomic. For example, incrementing a variable is not atomic.
func TestAtomacity(t *testing.T) {
	counter := 0
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()

	if counter != 1000 {
		t.Errorf("Expected counter to be 1000, got %d", counter)
	}
}

// Deadlocks and Livelocks
// Deadlock is a situation where two or more goroutines are waiting for each other to release resources, causing them to be blocked indefinitely.
// Livelock is a situation where two or more goroutines are actively trying to resolve a deadlock, but none of them can make progress.

func EatPasta(t *testing.T, name string, result chan<- string, cutlery ...*sync.Mutex) {
	for i, c := range cutlery {
		c.Lock()
		t.Logf("%s got %d item\n", name, i+1)
		time.Sleep(time.Microsecond)
		defer c.Unlock()
	}

	result <- name + " is done eating pasta"
}

func TestDeadlock(t *testing.T) {
	results := make(chan string)

	fork := &sync.Mutex{}
	spoon := &sync.Mutex{}

	go EatPasta(t, "Plato", results, fork, spoon)
	go EatPasta(t, "Socrates", results, spoon, fork)

	for i := 0; i < 2; i++ {
		select {
		case res := <-results:
			t.Log(res)
		case <-time.After(1 * time.Second):
			t.Fatal("Expected to receive result")
		}
	}
}

// Synchonization primitives, what do we have in Go?
// - Channels
// - Locks
// - Atomics
// What is the difference and when to use each of them?
// - if it's critical section, use Mutex or Atomics
// - if you're transferring ownership of data, use Channels
// - if you're trying to protect an internal state of a struct, use Mutex
// - if you're trying to coordinate multiple pieces of code, use Channels

// TODO: Would be nice to have a test for each of the synchronization primitives

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
