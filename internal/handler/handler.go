package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type URLService interface {
	CreateShortURL(longURL string, ttlMinutes *int) (string, error)
	Resolve(code string) (string, time.Time, bool, error)
}

type URLHandler struct {
	service URLService
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{service: service}
}

type shortenRequest struct {
	LongURL    string `json:"long_url" binding:"required,url"`
	TTLMinutes *int   `json:"ttl_minutes"` // pointer = optional
}

func (h *URLHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *URLHandler) Shorten(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	code, err := h.service.CreateShortURL(req.LongURL, req.TTLMinutes)
	if err != nil {
		log.Println("CREATE SHORT URL ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"short_url": "http://localhost:8080/" + code,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	url, expiresAt, ok, err := h.service.Resolve(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if !ok || time.Now().UTC().After(expiresAt) {
		c.JSON(http.StatusNotFound, gin.H{"error": "short url expired or not found"})
		return
	}

	c.Redirect(http.StatusFound, url)
}
