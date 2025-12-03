package helpers

import (
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.String(401, "Unauthorized: missing or invalid token")
			ctx.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]
		userId, role, err := ValidateJWTToken(tokenString)
		if err != nil {
			ctx.String(401, "Unauthorized: invalid token")
			ctx.Abort()
			return
		}
		ctx.Set("user_id", userId)
		ctx.Set("role", role)
		ctx.Next()
	}
}
