package goroutine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ProducerConsumer pattern
func ProducerConsumer() {
	fmt.Println("\n=== Producer-Consumer Pattern ===")
	
	buffer := make(chan int, 5)
	var wg sync.WaitGroup
	
	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(buffer)
		
		for i := 1; i <= 10; i++ {
			fmt.Printf("Producing: %d\n", i)
			buffer <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		for item := range buffer {
			fmt.Printf("Consuming: %d\n", item)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	
	wg.Wait()
}

// Semaphore pattern untuk membatasi concurrent access
func SemaphorePattern() {
	fmt.Println("\n=== Semaphore Pattern ===")
	
	// Semaphore dengan kapasitas 2
	semaphore := make(chan struct{}, 2)
	var wg sync.WaitGroup
	
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore
			
			fmt.Printf("Worker %d acquired semaphore\n", id)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("Worker %d releasing semaphore\n", id)
		}(i)
	}
	
	wg.Wait()
}

// Broadcast pattern menggunakan close channel
func BroadcastPattern() {
	fmt.Println("\n=== Broadcast Pattern ===")
	
	broadcast := make(chan struct{})
	var wg sync.WaitGroup
	
	// Start multiple listeners
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			fmt.Printf("Listener %d waiting for broadcast\n", id)
			<-broadcast
			fmt.Printf("Listener %d received broadcast signal\n", id)
		}(i)
	}
	
	// Wait a bit then broadcast
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Broadcasting signal...")
	close(broadcast) // This broadcasts to all listeners
	
	wg.Wait()
}

// Timeout pattern
func TimeoutPattern() {
	fmt.Println("\n=== Timeout Pattern ===")
	
	result := make(chan string)
	
	go func() {
		// Simulate long running task
		time.Sleep(2 * time.Second)
		result <- "Task completed"
	}()
	
	select {
	case res := <-result:
		fmt.Println("Success:", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Task timed out")
	}
}

// Heartbeat pattern
func HeartbeatPattern() {
	fmt.Println("\n=== Heartbeat Pattern ===")
	
	heartbeat := make(chan struct{})
	done := make(chan struct{})
	
	// Worker with heartbeat
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				select {
				case heartbeat <- struct{}{}:
					fmt.Println("💓 Heartbeat sent")
				default:
					// Non-blocking send
				}
			case <-done:
				fmt.Println("Worker stopping")
				return
			}
		}
	}()
	
	// Monitor heartbeat
	go func() {
		for {
			select {
			case <-heartbeat:
				fmt.Println("📡 Heartbeat received")
			case <-time.After(1 * time.Second):
				fmt.Println("⚠️  No heartbeat received - worker might be dead")
				return
			}
		}
	}()
	
	time.Sleep(3 * time.Second)
	close(done)
	time.Sleep(100 * time.Millisecond) // Let goroutines finish
}

// Retry pattern dengan exponential backoff
func RetryPattern() {
	fmt.Println("\n=== Retry Pattern ===")
	
	maxRetries := 3
	baseDelay := 100 * time.Millisecond
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := simulateFailingTask(); err != nil {
			if attempt == maxRetries {
				fmt.Printf("❌ Final attempt failed: %v\n", err)
				break
			}
			
			delay := time.Duration(attempt) * baseDelay
			fmt.Printf("🔄 Attempt %d failed: %v. Retrying in %v\n", attempt, err, delay)
			time.Sleep(delay)
		} else {
			fmt.Printf("✅ Task succeeded on attempt %d\n", attempt)
			break
		}
	}
}

func simulateFailingTask() error {
	if rand.Float32() < 0.7 { // 70% chance of failure
		return fmt.Errorf("simulated failure")
	}
	return nil
}

// Circuit breaker pattern
type CircuitBreaker struct {
	maxFailures int
	failures    int
	lastFailure time.Time
	timeout     time.Duration
	mu          sync.Mutex
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Check if circuit is open
	if cb.failures >= cb.maxFailures {
		if time.Since(cb.lastFailure) < cb.timeout {
			return fmt.Errorf("circuit breaker is open")
		}
		// Reset after timeout
		cb.failures = 0
	}
	
	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		return err
	}
	
	cb.failures = 0
	return nil
}

func CircuitBreakerPattern() {
	fmt.Println("\n=== Circuit Breaker Pattern ===")
	
	cb := NewCircuitBreaker(2, 2*time.Second)
	
	for i := 1; i <= 8; i++ {
		err := cb.Call(simulateFailingTask)
		if err != nil {
			fmt.Printf("Call %d failed: %v\n", i, err)
		} else {
			fmt.Printf("Call %d succeeded\n", i)
		}
		time.Sleep(300 * time.Millisecond)
	}
}