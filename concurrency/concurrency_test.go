package concurrency

import (
	"bytes"
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

// Channels patterns

// Fan-in Fan-out
// Fan-in is a pattern where multiple goroutines write to the same channel.
// This is useful when you have multiple goroutines producing results that need to be collected by a single goroutine.
// Fan-out is a pattern where a single goroutine reads from multiple channels.
// This is useful when you have a single goroutine that distributes work to multiple worker goroutines.

func Producer(n int, fanOut chan<- int) {
	// Here we should send n numbers to the fanOut channel
}

func Consumer(in <-chan int, out chan<- int) {
	// Here we should receive numbers from the in channel and send the result to the out channel
}

func TestFanInFanOut(t *testing.T) {
	expectedNums := 10
	work := make(chan int)
	results := make(chan int)

	for i := 0; i < 3; i++ {
		go Consumer(work, results)
	}

	go Producer(expectedNums, work)

	for i := 0; i < expectedNums; i++ {
		select {
		case res := <-results:
			t.Log(res)
		case <-time.After(1 * time.Second):
			t.Fatal("Expected to receive result")
		}
	}
}

// Channels of channels is a common pattern in Go to implement a producer-consumer model.
// In this pattern, we have a channel that is used to send requests to a producer goroutine.
// The producer goroutine processes the requests and sends the results back to the caller using a response channel.
type NumberIterator struct {
	requests chan chan<- int
	ctx      context.Context
}

func NewNumberIterator(ctx context.Context) *NumberIterator {
	return &NumberIterator{
		requests: make(chan chan<- int),
		ctx:      ctx,
	}
}

func (ni *NumberIterator) Next(ctx context.Context) (int, error) {
	return 0, nil
}

func (ni *NumberIterator) Run() {
	counter := 0

	for {
		select {
		case respChan := <-ni.requests:
			// Simulate some work
			time.Sleep(1 * time.Millisecond)
			counter++
			respChan <- counter
		case <-ni.ctx.Done():
			return
		}
	}
}

func TestChanOfChan(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ni := NewNumberIterator(ctx)

	go ni.Run()

	for i := 0; i < 2; i++ {
		num, err := ni.Next(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if num != i+1 {
			t.Errorf("Expected next number to be %d", i+1)
		}
	}
}

// Defalut case in select statement is used to handle non-blocking channel operations.
// If there are no other cases ready to be executed, the default case will be executed.
func TestDefaultCase(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ni := NewNumberIterator(ctx)
	go ni.Run()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	num, err := ni.Next(canceledCtx)
	if err != context.Canceled {
		t.Fatalf("Expected error to be %v, got %v", context.Canceled, err)
	}

	if num != 0 {
		t.Fatalf("Expected number to be 0, got %d", num)
	}

	num, err = ni.Next(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if num != 2 {
		t.Fatalf("Expected number to be 1, got %d", num)
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

// Sync.Pool is a synchronization primitive that is used to cache and reuse objects.
// It is useful for reducing memory allocations and improving performance.

func TestSyncPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			return bytes.NewBuffer([]byte{})
		},
	}

	buf := pool.Get().(*bytes.Buffer)
	buf.Write([]byte("hello"))

	if buf.String() != "hello" {
		t.Error("Expected data to be hello")
	}

	pool.Put(buf)

	buf = pool.Get().(*bytes.Buffer)
	if buf.String() != "" {
		t.Errorf("Expected data to be emptym got %s", buf.String())
	}
}

// Sync.Once is a synchronization primitive that guarantees that a function is executed only once.
// It is useful for initializing resources that are expensive to create or need to be shared across multiple goroutines.
// The Do method takes a function as an argument and ensures that the function is executed only once.

type RateLimiter struct {
	capacity int32
	counter  *atomic.Int32
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewRateLimiter(ctx context.Context, capacity int32) *RateLimiter {
	ctx, cancel := context.WithCancel(ctx)

	return &RateLimiter{
		capacity: capacity,
		counter:  &atomic.Int32{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (r *RateLimiter) Allow() bool {
	return r.counter.Add(1) <= r.capacity
}

func (r *RateLimiter) Close() {
	r.cancel()
}

func (r *RateLimiter) bucketRefiller() {
	t := time.NewTicker(1 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-t.C:
			r.counter.Store(0)
		}
	}
}

func TestSyncOnce(t *testing.T) {
	rl := NewRateLimiter(context.Background(), 3)
	defer rl.Close()

	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Error("Expected to allow access")
		}
	}

	if rl.Allow() {
		t.Error("Expected to deny access")
	}

	time.Sleep(1 * time.Millisecond)

	if !rl.Allow() {
		t.Error("Expected to allow access")
	}
}
