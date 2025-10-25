package routes

import (
	"user-authentication/database"
	"user-authentication/helpers"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func userLogin(ctx *gin.Context) {
	var body helpers.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.String(400, "Invalid request %s", err)
		return
	}
	email := body.Email
	// password := body.Password

	db, err := database.GetDatabaseClient()
	if err != nil {
		ctx.String(500, "Database connection error: %s, %s", err, db)
		return
	}
	var storedPassword, id string
	err = db.QueryRow("SELECT password, id FROM users WHERE email=$1", email).Scan(&storedPassword, &id)
	if err != nil {
		ctx.String(401, "Invalid username or password")
		return
	}
	token, err := helpers.CreateJWTToken(id)
	if err != nil {
		ctx.String(500, "Error creating JWT token: %s", err)
	}
	ctx.SetCookie("user-token", token, 3600, "/", "localhost", false, true)
	ctx.String(200, "Login successful")
}

func userSignUp(ctx *gin.Context) {
	var body helpers.SignupRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.String(400, "Invalid Request Parameters: %s", err)
		return
	}
	db, err := database.GetDatabaseClient()
	if err != nil {
		ctx.String(500, "Database connection error: %s, %s", err, db)
		return
	}
	_, err = db.Exec("INSERT INTO users (id, first_name, last_name, email, user_type, image, password, version, created_at, updated_at, permission, is_verified) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, 1, NOW(), NOW(), 'user', false)",
		body.FirstName, body.LastName, body.Email, body.Type, "{}", body.Password)
	if err != nil {
		ctx.String(500, "Error creating user: %s", err)
		return
	}
	ctx.String(200, "user created successfully")
}

func loginWithGoogle(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "google" {
		ctx.String(400, "Unsupported provider: %s", provider)
		return
	}
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func callBackHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	if provider != "google" {
		ctx.String(400, "Unsupported provider: %s", provider)
		return
	}
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	userId, exists, err := helpers.NewUserCheck(user.Email)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	if !exists {
		userId, err = helpers.CreateUser(user.Email, user.FirstName, user.LastName, "", "user")
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
	}

	token, err := helpers.CreateJWTToken(userId)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}
	ctx.SetCookie("user-token", token, 3600, "/", "localhost", false, true) // how do you share the token with frontend?
	ctx.String(200, "Authentication successful. Token: %s", token)
}
