package proxy

import (
	"net/http/httputil"
	"net/url"
	"user-authentication/helpers"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(ctx *gin.Context) {
		// 1. Check Access Token
		accessToken, err := ctx.Cookie("access_token")
		var userId, role string
		var tokenValid bool

		if err == nil && accessToken != "" {
			userId, role, err = helpers.ValidateJWTToken(accessToken)
			if err == nil {
				tokenValid = true
			}
		}

		// 2. If Access Token Invalid/Expired, Try Refresh
		if !tokenValid {
			refreshToken, err := ctx.Cookie("refresh_token")
			if err != nil || refreshToken == "" {
				ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: No valid tokens"})
				return
			}

			userId, role, err = helpers.ValidateRefreshToken(refreshToken)
			if err != nil {
				// Refresh token invalid/revoked -> Logout
				ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
				ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
				ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Session expired"})
				return
			}

			// Refresh Token Valid -> Issue New Pair
			newAccessToken, newRefreshToken, err := helpers.CreateTokenPair(userId, role)
			if err != nil {
				ctx.AbortWithStatusJSON(500, gin.H{"error": "Failed to refresh token"})
				return
			}

			// Set new cookies
			ctx.SetCookie("access_token", newAccessToken, 900, "/", "localhost", false, true)
			ctx.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", false, true)

			// Update accessToken variable for the downstream request
			accessToken = newAccessToken
		}

		// 3. Forward Request
		// Add user info headers for downstream service
		ctx.Request.Header.Set("X-User-ID", userId)
		ctx.Request.Header.Set("X-User-Role", role)
		// Ensure Authorization header is set with valid access token
		ctx.Request.Header.Set("Authorization", "Bearer "+accessToken)

		// Update host to match target
		ctx.Request.Host = targetURL.Host

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
