package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	user "vintage-server/internal/service/account" // Sesuaikan path
	"vintage-server/pkg/config"
)

func main() {
	// 1. Muat Konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 2. Koneksi Database menggunakan config
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// 2. Merakit semua lapisan (Wiring)
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg.JWTSecretKey) // Asumsi service.go sudah dibuat
	userHandler := user.NewHandler(userService)

	// 3. Setup Router Gin
	router := gin.Default()

	// 4. Daftarkan rute ke method di Handler
	// Ini adalah "API Contract" yang sesungguhnya
	api := router.Group("/api/v1") // Grup rute untuk versioning
	{
		account := api.Group("/account")
		{
			customer := account.Group("/customer")
			{
				customer.POST("/register", userHandler.RegisterCustomer)
				customer.POST("/login", userHandler.LoginCustomer)
			}
			admin := account.Group("/admin")
			{
				admin.POST("/login", userHandler.LoginAdmin)
			}
		}

	}

	// 5. Jalankan server
	log.Println("User Service running on port :8081")
	router.Run(":8081")
}
