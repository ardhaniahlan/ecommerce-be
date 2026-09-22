package routes

import (
	"ecommerce-backend/config"
	"ecommerce-backend/controllers"
	"ecommerce-backend/middlewares"
	"ecommerce-backend/repositories"
	"ecommerce-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(r *gin.Engine, redis *redis.Client) {
	productRepo := repositories.NewProductRepository(config.DB)
	productService := services.NewProductService(productRepo, redis)
	productController := controllers.NewProductController(productService)

	cartRepo := repositories.NewCartRepository(config.DB)
	cartService := services.NewCartService(cartRepo)
	cartController := controllers.NewCartController(cartService)

	orderRepo := repositories.NewOrderRepository(config.DB)
	orderService := services.NewOrderService(orderRepo)
	orderController := controllers.NewOrderController(orderService)

	userRepo := repositories.NewUserRepository(config.DB)
	service := services.NewUserService(userRepo)
	userController := controllers.NewUserController(service)

	voucherRepo := repositories.NewVoucherRepository(config.DB)
	voucherService := services.NewVoucherService(voucherRepo)
	voucherController := controllers.NewVoucherController(voucherService)

	bannerRepo := repositories.NewBannerRepository(config.DB)
	bannerService := services.NewBannerService(bannerRepo, redis)
	bannerController := controllers.NewBannerController(bannerService)

	

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		users := api.Group("/users", middlewares.RequireAuth())
		{
			users.GET("/profile", userController.GetProfile)
			users.PUT("/profile", userController.UpdateProfile)
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

		userOrders := api.Group("/orders", middlewares.RequireAuth())
		{
			userOrders.POST("/checkout", orderController.Checkout)
			userOrders.GET("/history", orderController.GetUserHistory)
			userOrders.PUT("/:id/complete", orderController.CompleteOrder)
		}

		api.POST("/payments/webhook", orderController.Webhook)

		adminGroup := api.Group("/admin", middlewares.RequireAuth(), middlewares.RequireAdmin())
		{
			adminGroup.GET("/stats", orderController.GetStats)

			adminOrders := adminGroup.Group("/orders")
			{
				adminOrders.GET("/", orderController.GetAllOrdersAdmin)
				adminOrders.PUT("/:id/tracking", orderController.InputTracking)
			}

			voucherRoutes := adminGroup.Group("/vouchers")
			{
				voucherRoutes.POST("/", voucherController.CreateVoucher)
				voucherRoutes.GET("/", voucherController.GetAllVoucher)
			}

			adminGroup.POST("/banners", bannerController.CreateBanner)
			adminGroup.DELETE("/banners/:id", bannerController.DeleteBanner)
		}
	}
	api.GET("/banners", bannerController.GetActiveBanners)
}
