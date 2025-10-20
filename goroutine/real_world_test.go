package goroutine

import (
	"testing"
)

func TestWebScraper(t *testing.T) {
	urls := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404", 
		"https://httpbin.org/delay/1",
	}
	
	results := WebScraper(urls)
	
	if len(results) != len(urls) {
		t.Errorf("Expected %d results, got %d", len(urls), len(results))
	}
}

func TestImageProcessor(t *testing.T) {
	ImageProcessor()
}

func TestAPIAggregator(t *testing.T) {
	APIAggregator()
}

func TestLogProcessor(t *testing.T) {
	LogProcessor()
}

func TestChatServer(t *testing.T) {
	ChatServer()
}