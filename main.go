package main

import (
	"fmt"
	"os"
	"user-authentication/database"
	"user-authentication/routes"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func startServer() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
	}
	_, err = database.GetDatabaseClient()
	if err != nil {
		fmt.Printf("Cannot connect to database: %s", err)
		return
	}
	key := "super"
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientCallback := "http://localhost:8009/auth/google/callback"
	maxAge := 86400 * 30
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // or false for dev
	gothic.Store = store
	goth.ClearProviders()
	goth.UseProviders(google.New(clientID, clientSecret, clientCallback))
	routes.SetupRoutes(router)
	router.Run(":8009")
}

func main() {
	startServer()
}
