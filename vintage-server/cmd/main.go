package main

import (
	"vintage-server/config"
	"vintage-server/models"
	"vintage-server/routes"
)

func main() {
	config.ConnectDatabase()

	// Auto migrate schema
	config.DB.AutoMigrate(&models.User{})

	// Setup routes
	r := routes.SetupRouter()
	r.Run(":8080")
}
