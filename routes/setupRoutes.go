package routes

import (
	"user-authentication/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
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

	router.POST("/login", userLogin)
	router.POST("/signup", userSignUp)
}
