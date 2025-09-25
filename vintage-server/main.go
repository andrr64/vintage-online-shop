package main

import (
	"log"
	"vintage-server/config"
	_ "vintage-server/docs"
	"vintage-server/routes"
	"vintage-server/shared"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect to database
	shared.ConnectDatabase()

	// Setup routes
	r := routes.SetupRouter()

	// Jalankan server di port 8080 (atau dari env)
	port := ":" + config.AppConfig.ServerPort // misal ambil dari env

	log.Printf("Starting server at %s...\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
