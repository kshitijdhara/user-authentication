package main

import (
	"fmt"
	"user-authentication/database"
	"user-authentication/routes"

	"github.com/gin-gonic/gin"
)

func startServer() {
	router := gin.Default()
	_, err := database.GetDatabaseClient()
	if err != nil {
		fmt.Printf("Cannot connect to database: %s", err)
		return
	}
	routes.SetupRoutes(router)
	router.Run(":8009")
}

func main() {
	startServer()
}
