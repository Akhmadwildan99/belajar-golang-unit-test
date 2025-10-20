package goroutine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 1. WEB SCRAPER - Scraping multiple URLs concurrently
type ScrapResult struct {
	URL    string
	Status int
	Size   int
	Error  error
}

func WebScraper(urls []string) []ScrapResult {
	fmt.Println("=== Web Scraper ===")
	
	results := make(chan ScrapResult, len(urls))
	var wg sync.WaitGroup
	
	// Scrape each URL concurrently
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			
			resp, err := http.Get(u)
			if err != nil {
				results <- ScrapResult{URL: u, Error: err}
				return
			}
			defer resp.Body.Close()
			
			body, _ := io.ReadAll(resp.Body)
			results <- ScrapResult{
				URL:    u,
				Status: resp.StatusCode,
				Size:   len(body),
			}
		}(url)
	}
	
	// Close results channel when all done
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Collect results
	var allResults []ScrapResult
	for result := range results {
		allResults = append(allResults, result)
		if result.Error != nil {
			fmt.Printf("❌ %s: %v\n", result.URL, result.Error)
		} else {
			fmt.Printf("✅ %s: %d bytes (status: %d)\n", result.URL, result.Size, result.Status)
		}
	}
	
	return allResults
}

// 2. IMAGE PROCESSOR - Process multiple images concurrently
type ImageJob struct {
	ID       int
	Filename string
	Width    int
	Height   int
}

type ProcessedImage struct {
	ID       int
	Filename string
	Success  bool
	Duration time.Duration
}

func ImageProcessor() {
	fmt.Println("\n=== Image Processor ===")
	
	jobs := []ImageJob{
		{1, "photo1.jpg", 1920, 1080},
		{2, "photo2.jpg", 1280, 720},
		{3, "photo3.jpg", 800, 600},
		{4, "photo4.jpg", 1024, 768},
		{5, "photo5.jpg", 640, 480},
	}
	
	jobChan := make(chan ImageJob, len(jobs))
	resultChan := make(chan ProcessedImage, len(jobs))
	
	// Start 3 image processing workers
	numWorkers := 3
	var wg sync.WaitGroup
	
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go imageWorker(i, jobChan, resultChan, &wg)
	}
	
	// Send jobs
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)
	
	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	for result := range resultChan {
		status := "❌"
		if result.Success {
			status = "✅"
		}
		fmt.Printf("%s Image %d (%s) processed in %v\n", 
			status, result.ID, result.Filename, result.Duration)
	}
}

func imageWorker(id int, jobs <-chan ImageJob, results chan<- ProcessedImage, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for job := range jobs {
		start := time.Now()
		fmt.Printf("Worker %d processing image %d (%s)\n", id, job.ID, job.Filename)
		
		// Simulate image processing (resize, compress, etc.)
		processingTime := time.Duration(job.Width*job.Height/100000) * time.Millisecond
		time.Sleep(processingTime)
		
		results <- ProcessedImage{
			ID:       job.ID,
			Filename: job.Filename,
			Success:  true,
			Duration: time.Since(start),
		}
	}
}

// 3. API AGGREGATOR - Fetch data from multiple APIs
type APIResponse struct {
	Service string
	Data    interface{}
	Error   error
	Latency time.Duration
}

func APIAggregator() {
	fmt.Println("\n=== API Aggregator ===")
	
	services := []string{"users", "products", "orders", "analytics"}
	responses := make(chan APIResponse, len(services))
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// Fetch from all APIs concurrently
	for _, service := range services {
		go fetchFromAPI(ctx, service, responses)
	}
	
	// Collect responses with timeout
	results := make(map[string]APIResponse)
	for i := 0; i < len(services); i++ {
		select {
		case resp := <-responses:
			results[resp.Service] = resp
			status := "✅"
			if resp.Error != nil {
				status = "❌"
			}
			fmt.Printf("%s %s API: %v (latency: %v)\n", 
				status, resp.Service, resp.Data, resp.Latency)
		case <-ctx.Done():
			fmt.Printf("⏰ Timeout waiting for API responses\n")
			return
		}
	}
}

func fetchFromAPI(ctx context.Context, service string, responses chan<- APIResponse) {
	start := time.Now()
	
	// Simulate API call with random delay
	select {
	case <-time.After(time.Duration(500+int(service[0])*10) * time.Millisecond):
		responses <- APIResponse{
			Service: service,
			Data:    fmt.Sprintf("%s_data_123", service),
			Latency: time.Since(start),
		}
	case <-ctx.Done():
		responses <- APIResponse{
			Service: service,
			Error:   ctx.Err(),
			Latency: time.Since(start),
		}
	}
}

// 4. LOG PROCESSOR - Process log files concurrently
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Service   string
}

func LogProcessor() {
	fmt.Println("\n=== Log Processor ===")
	
	// Simulate log entries
	logs := []LogEntry{
		{time.Now(), "ERROR", "Database connection failed", "api"},
		{time.Now(), "INFO", "User logged in", "auth"},
		{time.Now(), "WARN", "High memory usage", "worker"},
		{time.Now(), "ERROR", "Payment failed", "payment"},
		{time.Now(), "INFO", "Cache cleared", "cache"},
	}
	
	logChan := make(chan LogEntry, len(logs))
	errorChan := make(chan LogEntry, len(logs))
	
	// Start log processors
	var wg sync.WaitGroup
	
	// Error log processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for errorLog := range errorChan {
			fmt.Printf("🚨 ALERT: %s [%s] %s\n", 
				errorLog.Service, errorLog.Level, errorLog.Message)
			// Send to monitoring system, email, Slack, etc.
		}
	}()
	
	// General log processor
	wg.Add(1)
	go func() {
		defer wg.Done()
		for log := range logChan {
			if log.Level == "ERROR" {
				errorChan <- log
			}
			fmt.Printf("📝 %s [%s] %s: %s\n", 
				log.Timestamp.Format("15:04:05"), log.Level, log.Service, log.Message)
		}
		close(errorChan)
	}()
	
	// Send logs
	for _, log := range logs {
		logChan <- log
	}
	close(logChan)
	
	wg.Wait()
}

// 5. REAL-TIME CHAT SERVER
type ChatMessage struct {
	UserID  string
	RoomID  string
	Message string
	Time    time.Time
}

type ChatRoom struct {
	ID      string
	clients map[string]chan ChatMessage
	mu      sync.RWMutex
}

func (cr *ChatRoom) AddClient(userID string) chan ChatMessage {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	
	ch := make(chan ChatMessage, 10)
	cr.clients[userID] = ch
	return ch
}

func (cr *ChatRoom) RemoveClient(userID string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	
	if ch, exists := cr.clients[userID]; exists {
		close(ch)
		delete(cr.clients, userID)
	}
}

func (cr *ChatRoom) Broadcast(msg ChatMessage) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	
	for userID, ch := range cr.clients {
		if userID != msg.UserID { // Don't send to sender
			select {
			case ch <- msg:
			default:
				// Channel full, skip
			}
		}
	}
}

func ChatServer() {
	fmt.Println("\n=== Chat Server ===")
	
	room := &ChatRoom{
		ID:      "general",
		clients: make(map[string]chan ChatMessage),
	}
	
	// Simulate 3 users joining
	users := []string{"alice", "bob", "charlie"}
	var wg sync.WaitGroup
	
	for _, user := range users {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()
			
			// Join room
			msgChan := room.AddClient(userID)
			defer room.RemoveClient(userID)
			
			fmt.Printf("👤 %s joined the room\n", userID)
			
			// Listen for messages
			go func() {
				for msg := range msgChan {
					fmt.Printf("📨 %s received: [%s] %s\n", 
						userID, msg.UserID, msg.Message)
				}
			}()
			
			// Send some messages
			messages := []string{
				fmt.Sprintf("Hello from %s!", userID),
				fmt.Sprintf("%s is typing...", userID),
			}
			
			for _, text := range messages {
				msg := ChatMessage{
					UserID:  userID,
					RoomID:  room.ID,
					Message: text,
					Time:    time.Now(),
				}
				room.Broadcast(msg)
				time.Sleep(500 * time.Millisecond)
			}
			
			time.Sleep(2 * time.Second) // Stay in room
		}(user)
	}
	
	wg.Wait()
}