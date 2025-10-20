# Belajar Goroutine - Panduan Lengkap

## 📚 Daftar Isi

### 1. Basic Goroutine (`basic_goroutine.go`)
- **BasicGoroutine()** - Penggunaan dasar goroutine
- **GoroutineWithWaitGroup()** - Sinkronisasi dengan WaitGroup
- **ChannelCommunication()** - Komunikasi antar goroutine
- **BufferedChannel()** - Channel dengan buffer
- **SelectStatement()** - Penggunaan select untuk multiple channel

### 2. Advanced Patterns (`advanced_goroutine.go`)
- **WorkerPool()** - Pattern worker pool untuk concurrent processing
- **ContextCancellation()** - Pembatalan goroutine dengan context
- **PipelinePattern()** - Pipeline untuk data processing
- **FanOutFanIn()** - Distribusi dan agregasi data
- **RateLimiter()** - Pembatasan rate request

### 3. Complex Examples (`examples.go`)
- **ProducerConsumer()** - Pattern producer-consumer
- **SemaphorePattern()** - Pembatasan concurrent access
- **BroadcastPattern()** - Broadcasting signal ke multiple goroutine
- **TimeoutPattern()** - Handling timeout
- **HeartbeatPattern()** - Monitoring goroutine health
- **RetryPattern()** - Retry dengan exponential backoff
- **CircuitBreakerPattern()** - Circuit breaker untuk fault tolerance

## 🚀 Cara Menjalankan

### Menjalankan Semua Contoh
```bash
cd goroutine
go run main.go
```

### Menjalankan Test
```bash
# Test semua
go test -v

# Test dengan race detection
go test -race -v

# Benchmark
go test -bench=. -v
```

### Menjalankan Contoh Spesifik
```go
package main

import "github.com/belajar-golang-unit-test/goroutine"

func main() {
    goroutine.BasicGoroutine()
    goroutine.WorkerPool()
    // dst...
}
```

## 📖 Konsep Penting

### 1. Goroutine
- Lightweight thread yang dikelola Go runtime
- Dimulai dengan keyword `go`
- Sangat efisien (hanya ~2KB memory overhead)

### 2. Channel
- Cara komunikasi antar goroutine
- **Unbuffered**: Synchronous communication
- **Buffered**: Asynchronous communication dengan kapasitas tertentu

### 3. Select Statement
- Memungkinkan goroutine menunggu multiple channel operations
- Non-blocking dengan `default` case
- Timeout dengan `time.After()`

### 4. Synchronization
- **WaitGroup**: Menunggu sekelompok goroutine selesai
- **Mutex**: Mutual exclusion untuk shared resources
- **Context**: Cancellation dan timeout

## 🎯 Pattern yang Dipelajari

### Worker Pool
```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// Start workers
for w := 1; w <= numWorkers; w++ {
    go worker(w, jobs, results)
}
```

### Pipeline
```go
numbers := generate(1, 2, 3, 4, 5)
squares := square(numbers)
doubled := double(squares)
```

### Fan-Out Fan-In
```go
// Fan-out
c1 := process(input)
c2 := process(input)

// Fan-in
output := merge(c1, c2)
```

## ⚠️ Best Practices

1. **Selalu tutup channel** setelah selesai mengirim data
2. **Gunakan WaitGroup** untuk sinkronisasi
3. **Handle race conditions** dengan mutex atau channel
4. **Gunakan context** untuk cancellation
5. **Jangan buat terlalu banyak goroutine** - gunakan worker pool
6. **Test dengan `-race` flag** untuk detect race conditions

## 🔧 Debugging Tips

### Race Condition Detection
```bash
go test -race -v
```

### Goroutine Leak Detection
```bash
go test -v -count=1 ./...
```

### Memory Profiling
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## 📊 Performance Notes

- Goroutine creation: ~1000ns
- Channel communication: ~100ns
- Context switching: ~200ns
- Mutex lock/unlock: ~20ns

## 🎓 Latihan

1. Implementasikan web crawler concurrent
2. Buat rate limiter dengan sliding window
3. Implementasikan distributed work queue
4. Buat monitoring system dengan heartbeat
5. Implementasikan graceful shutdown

## 📚 Referensi

- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)
- [Advanced Go Concurrency Patterns](https://talks.golang.org/2013/advconc.slide)
- [Go Memory Model](https://golang.org/ref/mem)