package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type URLStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		data: make(map[string]string),
	}
}

func (s *URLStore) Save(shortCode, longURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[shortCode] = longURL
}

func (s *URLStore) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, ok := s.data[shortCode]
	return longURL, ok
}

func main() {

	store := NewURLStore()
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		store.Save("abc123", "https://example.com")
		c.JSON(http.StatusOK, gin.H{"message": "saved"})
	})

	r.Run(":8080")
}
