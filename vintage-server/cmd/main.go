package main

import (
	"log"
	"vintage-server/config"
	_ "vintage-server/docs"
	"vintage-server/shared"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Connect to database
	shared.ConnectDatabase()

	log.Println("Project initialized successfully")
	// Di sini nanti bisa start Gin server
}
