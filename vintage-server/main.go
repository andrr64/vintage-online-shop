package main

import (
	"vintage-server/config"
	_ "vintage-server/docs"
	"vintage-server/models"
	"vintage-server/routes"
)

// @title 	Vintage Server
// @version 1.0
// @basepath /api/v1

func main() {
	config.ConnectDatabase()

	// Auto migrate schema
	config.DB.AutoMigrate(&models.User{})

	// Setup routes
	r := routes.SetupRouter()

	r.Run(":8080")
}
