package routes

import (
	"user-authentication/database"
	"user-authentication/helpers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/login", userLogin)
	router.POST("/signup", userSignUp)
	router.GET("/auth/:provider", loginWithGoogle)
	router.GET("/auth/:provider/callback", callBackHandler)

	apirouter := router.Group("/api", helpers.AuthMiddleware())
	{
		apirouter.GET("/", func(ctx *gin.Context) {
			db, err := database.GetDatabaseClient()
			if err != nil {
				ctx.JSON(500, gin.H{"db": "disconnected", "error": err.Error()})
				return
			}
			if err := db.Ping(); err != nil {
				ctx.JSON(500, gin.H{"db": "ping_failed", "error": err.Error()})
				return
			}
			ctx.JSON(200, gin.H{"db": "ok"})
		})
	}
}
