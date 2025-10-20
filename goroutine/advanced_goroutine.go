package goroutine

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WorkerPool implementasi worker pool pattern
func WorkerPool() {
	fmt.Println("\n=== Worker Pool Pattern ===")
	
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	
	// Start workers
	var wg sync.WaitGroup
	numWorkers := 3
	
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	
	// Send jobs
	numJobs := 5
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second)
		results <- job * 2
	}
}

// ContextCancellation menggunakan context untuk pembatalan
func ContextCancellation() {
	fmt.Println("\n=== Context Cancellation ===")
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	go func() {
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("Task completed")
		case <-ctx.Done():
			fmt.Println("Task cancelled:", ctx.Err())
		}
	}()
	
	time.Sleep(3 * time.Second)
}

// PipelinePattern implementasi pipeline pattern
func PipelinePattern() {
	fmt.Println("\n=== Pipeline Pattern ===")
	
	numbers := generate(1, 2, 3, 4, 5)
	squares := square(numbers)
	doubled := double(squares)
	
	for result := range doubled {
		fmt.Printf("Result: %d\n", result)
	}
}

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func double(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * 2
		}
	}()
	return out
}

// FanOutFanIn pattern untuk distribusi dan agregasi
func FanOutFanIn() {
	fmt.Println("\n=== Fan-Out Fan-In Pattern ===")
	
	input := make(chan int)
	
	// Fan-out: distribute work to multiple goroutines
	c1 := process(input)
	c2 := process(input)
	c3 := process(input)
	
	// Fan-in: merge results
	output := merge(c1, c2, c3)
	
	// Send input
	go func() {
		defer close(input)
		for i := 1; i <= 6; i++ {
			input <- i
		}
	}()
	
	// Collect results
	for result := range output {
		fmt.Printf("Processed: %d\n", result)
	}
}

func process(input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for n := range input {
			time.Sleep(100 * time.Millisecond) // Simulate work
			output <- n * n
		}
	}()
	return output
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	output := make(chan int)
	
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			output <- n
		}
	}
	
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	
	go func() {
		wg.Wait()
		close(output)
	}()
	
	return output
}

// RateLimiter implementasi rate limiting
func RateLimiter() {
	fmt.Println("\n=== Rate Limiter ===")
	
	// Rate limiter: 2 requests per second
	limiter := time.Tick(500 * time.Millisecond)
	
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)
	
	for req := range requests {
		<-limiter // Wait for rate limiter
		fmt.Printf("Processing request %d at %s\n", req, time.Now().Format("15:04:05"))
	}
}