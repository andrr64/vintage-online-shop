package main

import (
	"vintage-server/config"
	_ "vintage-server/docs"
)

// @title 	Vintage Server
// @version 1.0
// @host localhost:8080
// @basepath /api/v1

func main() {
	config.ConnectDatabase()
}
