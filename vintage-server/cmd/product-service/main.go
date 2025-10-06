package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"log"
	"vintage-server/internal/database"
	"vintage-server/internal/service/product"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/config"
	"vintage-server/pkg/middleware"
	"vintage-server/pkg/uploader"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 1. Koneksi Database (tidak berubah)
	db, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// 2. Inisialisasi semua service eksternal (tidak berubah)
	cloudinaryService, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary service: %v", err)
	}
	authService := auth.NewJWTService(cfg.JWTSecretKey)

	// highlight-start
	// 3. RAKITAN ARSITEKTUR BARU

	// 3a. Buat instance Store, yang di dalamnya sudah ada Repository.
	// Kita tidak lagi membuat Repository secara langsung di main.
	productStore := product.NewStore(db)

	// 3b. Suntikkan (inject) Store dan service lain ke dalam Product Service.
	// NewService sekarang menerima Store, bukan Repository.
	productService := product.NewService(productStore, *authService, cloudinaryService)

	// 3c. Suntikkan Product Service ke dalam Handler (tidak berubah).
	productHandler := product.NewHandler(productService)
	// highlight-end

	// 4. Setup Router (tidak berubah)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		productGroup := api.Group("/product")
		{
			// Rute Publik
			productGroup.GET("/category", productHandler.ReadCategories)
			productGroup.GET("/brand", productHandler.ReadBrand)
			productGroup.GET("/condition", productHandler.ReadConditions)

			// Rute Terproteksi
			protected := productGroup.Group("/protected", middleware.AuthMiddleware(authService))
			{
				// Category
				protected.POST("/category", productHandler.CreateCategory)
				protected.PUT("/category/:id", productHandler.UpdateCategory) // Gunakan path param untuk konsistensi
				protected.DELETE("/category/:id", productHandler.DeleteCategory)

				// Brand
				protected.POST("/brand", productHandler.CreateBrand)
				protected.PUT("/brand/:id", productHandler.UpdateBrand)
				protected.DELETE("/brand/:id", productHandler.DeleteBrand)

				// Condition
				protected.POST("/condition", productHandler.CreateCondition)
				protected.PUT("/condition/:id", productHandler.UpdateCondition)
				protected.DELETE("/condition/:id", productHandler.DeleteCondition)

				// Product
				protected.POST("", productHandler.CreateProduct)
				// Tambahkan rute product lainnya di sini (GET, PUT, DELETE)
			}
		}
	}

	// 5. Jalankan Server (tidak berubah)
	log.Printf("Product Service running on port :%s", cfg.ProductServicePort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ProductServicePort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
