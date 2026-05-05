package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders додає стандартні захисні HTTP-заголовки до кожної відповіді.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Заборона вбудовування у фрейми (clickjacking)
		c.Header("X-Frame-Options", "DENY")
		// Заборона MIME-sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// Увімкнення вбудованого XSS-фільтра браузера
		c.Header("X-XSS-Protection", "1; mode=block")
		// HSTS — браузер завжди використовує HTTPS (1 рік)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Обмежуємо Referrer
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// Вимикаємо небезпечні функції браузера
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// Прибираємо заголовок, що розкриває технологію
		c.Header("X-Powered-By", "")
		c.Next()
	}
}

// CORSMiddleware дозволяє запити лише з авторизованих origins.
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowedSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowedSet[o] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" && allowedSet[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// --- Простий in-memory rate limiter ---

type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}
	// Фонова очистка застарілих записів кожні 5 хвилин
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	for key, e := range rl.entries {
		if now.After(e.windowEnd) {
			delete(rl.entries, key)
		}
	}
}

// Allow перевіряє, чи не перевищено ліміт для даного ключа (IP або email).
// Повертає true, якщо запит дозволено.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	e, exists := rl.entries[key]
	if !exists || now.After(e.windowEnd) {
		rl.entries[key] = &rateLimitEntry{count: 1, windowEnd: now.Add(rl.window)}
		return true
	}
	if e.count >= rl.limit {
		return false
	}
	e.count++
	return true
}

// RateLimitMiddleware обмежує кількість запитів з однієї IP-адреси.
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Забагато спроб. Спробуйте пізніше.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
