package goroutine

import (
	"sync"
	"testing"
	"time"
)

func TestBasicGoroutine(t *testing.T) {
	BasicGoroutine()
}

func TestGoroutineWithWaitGroup(t *testing.T) {
	GoroutineWithWaitGroup()
}

func TestChannelCommunication(t *testing.T) {
	ChannelCommunication()
}

func TestBufferedChannel(t *testing.T) {
	BufferedChannel()
}

func TestSelectStatement(t *testing.T) {
	SelectStatement()
}

func TestWorkerPool(t *testing.T) {
	WorkerPool()
}

func TestContextCancellation(t *testing.T) {
	ContextCancellation()
}

func TestPipelinePattern(t *testing.T) {
	PipelinePattern()
}

func TestFanOutFanIn(t *testing.T) {
	FanOutFanIn()
}

func TestRateLimiter(t *testing.T) {
	RateLimiter()
}

// Benchmark tests
func BenchmarkGoroutineCreation(b *testing.B) {
	var wg sync.WaitGroup
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Minimal work
		}()
	}
	wg.Wait()
}

func BenchmarkChannelCommunication(b *testing.B) {
	ch := make(chan int, 1000)
	
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
	}()
	
	b.ResetTimer()
	for range ch {
		// Receive all values
	}
}

// Test race condition detection
func TestRaceCondition(t *testing.T) {
	counter := 0
	var wg sync.WaitGroup
	
	// This will cause race condition if run with -race flag
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // Race condition here
		}()
	}
	
	wg.Wait()
	t.Logf("Counter value: %d (may vary due to race condition)", counter)
}

// Test with mutex to fix race condition
func TestWithMutex(t *testing.T) {
	counter := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	if counter != 100 {
		t.Errorf("Expected counter to be 100, got %d", counter)
	}
}

// Test channel closing
func TestChannelClosing(t *testing.T) {
	ch := make(chan int, 5)
	
	// Send some values
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()
	
	// Receive until channel is closed
	var received []int
	for value := range ch {
		received = append(received, value)
	}
	
	expected := []int{1, 2, 3}
	if len(received) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, received)
	}
}

// Test select with timeout
func TestSelectTimeout(t *testing.T) {
	ch := make(chan string)
	
	select {
	case msg := <-ch:
		t.Errorf("Unexpected message: %s", msg)
	case <-time.After(100 * time.Millisecond):
		// Expected timeout
		t.Log("Timeout occurred as expected")
	}
}