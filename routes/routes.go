package routes

import (
	"ecommerce-backend/config"
	"ecommerce-backend/controllers"
	"ecommerce-backend/middlewares"
	"ecommerce-backend/repositories"
	"ecommerce-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	productRepo := repositories.NewProductRepository(config.DB)
	productService := services.NewProductService(productRepo)
	productController := controllers.NewProductController(productService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		products := api.Group("/products")
		{
			products.GET("/", productController.GetAll)
			products.GET("/:id", productController.GetByID)
			
			products.POST("/", middlewares.RequireAuth(), middlewares.RequireAdmin(), productController.Create) 
			products.PUT("/:id", middlewares.RequireAuth(), middlewares.RequireAdmin(), productController.Update) 
			products.DELETE("/:id", middlewares.RequireAuth(), middlewares.RequireAdmin(), productController.Delete)
		}
	}
}