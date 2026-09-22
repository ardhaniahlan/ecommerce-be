package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(redisClient *redis.Client, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		cacheKey := fmt.Sprintf("ratelimit:%s", clientIP)

		ctx := context.Background()

		totalRequests, err := redisClient.Incr(ctx, cacheKey).Result()
		if err != nil {
			c.Next()
			return
		}

		if totalRequests == 1 {
			redisClient.Expire(ctx, cacheKey, window)
		}

		if totalRequests > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  429,
				"error":   "Too Many Requests",
				"message": "Anda melakukan terlalu banyak request. Silakan coba lagi beberapa saat.",
			})
			return
		}

		c.Next()
	}
}