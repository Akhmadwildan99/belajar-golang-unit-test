package goroutine

import (
	"fmt"
	"sync"
	"time"
)

// BasicGoroutine mendemonstrasikan penggunaan dasar goroutine
func BasicGoroutine() {
	fmt.Println("=== Basic Goroutine ===")
	
	// Goroutine sederhana
	go func() {
		fmt.Println("Hello from goroutine!")
	}()
	
	// Tunggu sebentar agar goroutine selesai
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main function finished")
}

// GoroutineWithWaitGroup menggunakan WaitGroup untuk sinkronisasi
func GoroutineWithWaitGroup() {
	fmt.Println("\n=== Goroutine with WaitGroup ===")
	
	var wg sync.WaitGroup
	
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d starting\n", id)
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("Worker %d finished\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Println("All workers finished")
}

// ChannelCommunication mendemonstrasikan komunikasi antar goroutine
func ChannelCommunication() {
	fmt.Println("\n=== Channel Communication ===")
	
	ch := make(chan string)
	
	go func() {
		ch <- "Hello from goroutine via channel!"
	}()
	
	message := <-ch
	fmt.Println("Received:", message)
}

// BufferedChannel menggunakan buffered channel
func BufferedChannel() {
	fmt.Println("\n=== Buffered Channel ===")
	
	ch := make(chan int, 3)
	
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("Sending %d\n", i)
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()
	
	for value := range ch {
		fmt.Printf("Received %d\n", value)
		time.Sleep(200 * time.Millisecond)
	}
}

// SelectStatement mendemonstrasikan penggunaan select
func SelectStatement() {
	fmt.Println("\n=== Select Statement ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch1 <- "Message from channel 1"
	}()
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- "Message from channel 2"
	}()
	
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		case <-time.After(300 * time.Millisecond):
			fmt.Println("Timeout!")
		}
	}
}