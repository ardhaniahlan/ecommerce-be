package main

import (
	"ecommerce-backend/config"
	"ecommerce-backend/routes"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel sistem bawaan OS.")
	}

	config.InitDB()

	r := gin.Default()

	routes.SetupRoutes(r)

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