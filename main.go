package main

import (
	"context"
	"crypto/tls"
	"ecommerce-backend/config"
	"ecommerce-backend/middlewares"
	"ecommerce-backend/routes"
	"time"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel sistem bawaan OS.")
	}

	config.InitDB()

	r := gin.Default()

	
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,

		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})
	
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic("Gagal connect ke Redis: " + err.Error())
	}
	
	r.Use(middlewares.RateLimiter(redisClient, 3, 1*time.Minute))

	routes.SetupRoutes(r, redisClient)

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Database terhubung, Backend E-Commerce Go siap digunakan!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
