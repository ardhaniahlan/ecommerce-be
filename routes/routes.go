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

	cartRepo := repositories.NewCartRepository(config.DB)
	cartService := services.NewCartService(cartRepo)
	cartController := controllers.NewCartController(cartService)

	orderRepo := repositories.NewOrderRepository(config.DB)
	orderService := services.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

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

		cart := api.Group("/cart", middlewares.RequireAuth())
		{
			cart.GET("/", cartController.GetCart)
			cart.POST("/", cartController.AddItem)
			cart.DELETE("/:id", cartController.RemoveItem)
			cart.PUT("/:id", cartController.UpdateQuantity)
		}

		orders := api.Group("/orders", middlewares.RequireAuth())
		{
			orders.POST("/checkout", orderController.Checkout)
		}
	}
}
