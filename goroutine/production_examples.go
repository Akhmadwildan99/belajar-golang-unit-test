package goroutine

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// 1. DATABASE CONNECTION POOL
type DBPool struct {
	connections chan *sql.DB
	maxConns    int
	mu          sync.Mutex
}

func NewDBPool(maxConns int) *DBPool {
	return &DBPool{
		connections: make(chan *sql.DB, maxConns),
		maxConns:    maxConns,
	}
}

func (p *DBPool) Get() *sql.DB {
	select {
	case conn := <-p.connections:
		return conn
	default:
		// Create new connection if pool not full
		return p.createConnection()
	}
}

func (p *DBPool) Put(conn *sql.DB) {
	select {
	case p.connections <- conn:
	default:
		// Pool full, close connection
		conn.Close()
	}
}

func (p *DBPool) createConnection() *sql.DB {
	// Simulate DB connection creation
	fmt.Println("Creating new DB connection")
	return &sql.DB{} // Mock connection
}

// 2. MICROSERVICE HEALTH CHECKER
type ServiceHealth struct {
	Name     string
	URL      string
	Status   string
	Latency  time.Duration
	LastCheck time.Time
}

type HealthChecker struct {
	services []string
	results  map[string]ServiceHealth
	mu       sync.RWMutex
}

func NewHealthChecker(services []string) *HealthChecker {
	return &HealthChecker{
		services: services,
		results:  make(map[string]ServiceHealth),
	}
}

func (hc *HealthChecker) StartMonitoring(ctx context.Context) {
	fmt.Println("=== Health Checker Started ===")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	// Initial check
	hc.checkAllServices()
	
	for {
		select {
		case <-ticker.C:
			hc.checkAllServices()
		case <-ctx.Done():
			fmt.Println("Health checker stopped")
			return
		}
	}
}

func (hc *HealthChecker) checkAllServices() {
	var wg sync.WaitGroup
	
	for _, service := range hc.services {
		wg.Add(1)
		go func(svc string) {
			defer wg.Done()
			hc.checkService(svc)
		}(service)
	}
	
	wg.Wait()
	hc.printStatus()
}

func (hc *HealthChecker) checkService(service string) {
	start := time.Now()
	
	// Simulate health check
	time.Sleep(time.Duration(100+len(service)*10) * time.Millisecond)
	
	status := "healthy"
	if len(service)%2 == 0 {
		status = "unhealthy"
	}
	
	hc.mu.Lock()
	hc.results[service] = ServiceHealth{
		Name:      service,
		Status:    status,
		Latency:   time.Since(start),
		LastCheck: time.Now(),
	}
	hc.mu.Unlock()
}

func (hc *HealthChecker) printStatus() {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	fmt.Println("\n--- Service Health Status ---")
	for _, health := range hc.results {
		icon := "✅"
		if health.Status != "healthy" {
			icon = "❌"
		}
		fmt.Printf("%s %s: %s (latency: %v)\n", 
			icon, health.Name, health.Status, health.Latency)
	}
}

// 3. BATCH JOB PROCESSOR
type BatchJob struct {
	ID       string
	Data     []interface{}
	Priority int
}

type JobResult struct {
	JobID    string
	Success  bool
	Duration time.Duration
	Error    error
}

type BatchProcessor struct {
	jobQueue    chan BatchJob
	resultQueue chan JobResult
	workers     int
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewBatchProcessor(workers int, queueSize int) *BatchProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &BatchProcessor{
		jobQueue:    make(chan BatchJob, queueSize),
		resultQueue: make(chan JobResult, queueSize),
		workers:     workers,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (bp *BatchProcessor) Start() {
	fmt.Printf("=== Starting Batch Processor with %d workers ===\n", bp.workers)
	
	var wg sync.WaitGroup
	
	// Start workers
	for i := 1; i <= bp.workers; i++ {
		wg.Add(1)
		go bp.worker(i, &wg)
	}
	
	// Start result collector
	go bp.collectResults()
	
	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(bp.resultQueue)
	}()
}

func (bp *BatchProcessor) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for {
		select {
		case job := <-bp.jobQueue:
			start := time.Now()
			fmt.Printf("Worker %d processing job %s\n", id, job.ID)
			
			// Simulate job processing
			processingTime := time.Duration(len(job.Data)*100) * time.Millisecond
			time.Sleep(processingTime)
			
			bp.resultQueue <- JobResult{
				JobID:    job.ID,
				Success:  true,
				Duration: time.Since(start),
			}
			
		case <-bp.ctx.Done():
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}

func (bp *BatchProcessor) collectResults() {
	for result := range bp.resultQueue {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		fmt.Printf("%s Job %s completed in %v\n", 
			status, result.JobID, result.Duration)
	}
}

func (bp *BatchProcessor) SubmitJob(job BatchJob) {
	select {
	case bp.jobQueue <- job:
		fmt.Printf("📝 Job %s queued\n", job.ID)
	default:
		fmt.Printf("⚠️ Job queue full, job %s rejected\n", job.ID)
	}
}

func (bp *BatchProcessor) Shutdown() {
	fmt.Println("🛑 Shutting down batch processor...")
	close(bp.jobQueue)
	bp.cancel()
}

// 4. CACHE WITH TTL
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type TTLCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewTTLCache(ttl time.Duration) *TTLCache {
	cache := &TTLCache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

func (c *TTLCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	
	return item.Value, true
}

func (c *TTLCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Demo functions
func DemoDBPool() {
	fmt.Println("=== Database Pool Demo ===")
	
	pool := NewDBPool(3)
	var wg sync.WaitGroup
	
	// Simulate 5 concurrent database operations
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			conn := pool.Get()
			fmt.Printf("Worker %d got DB connection\n", id)
			
			// Simulate database work
			time.Sleep(time.Duration(id*100) * time.Millisecond)
			
			pool.Put(conn)
			fmt.Printf("Worker %d returned DB connection\n", id)
		}(i)
	}
	
	wg.Wait()
}

func DemoHealthChecker() {
	services := []string{"api-gateway", "user-service", "payment-service", "notification-service"}
	
	checker := NewHealthChecker(services)
	
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	checker.StartMonitoring(ctx)
}

func DemoBatchProcessor() {
	processor := NewBatchProcessor(3, 10)
	processor.Start()
	
	// Submit some jobs
	jobs := []BatchJob{
		{"job-1", []interface{}{1, 2, 3}, 1},
		{"job-2", []interface{}{4, 5, 6, 7}, 2},
		{"job-3", []interface{}{8, 9}, 1},
		{"job-4", []interface{}{10, 11, 12, 13, 14}, 3},
	}
	
	for _, job := range jobs {
		processor.SubmitJob(job)
	}
	
	time.Sleep(3 * time.Second)
	processor.Shutdown()
	time.Sleep(1 * time.Second)
}

func DemoTTLCache() {
	fmt.Println("=== TTL Cache Demo ===")
	
	cache := NewTTLCache(2 * time.Second)
	
	// Set some values
	cache.Set("user:123", "John Doe")
	cache.Set("session:abc", "active")
	
	// Get values immediately
	if val, ok := cache.Get("user:123"); ok {
		fmt.Printf("Found: %v\n", val)
	}
	
	// Wait for expiration
	time.Sleep(3 * time.Second)
	
	if _, ok := cache.Get("user:123"); !ok {
		fmt.Println("Value expired and removed")
	}
}