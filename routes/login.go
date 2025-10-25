package routes

import (
	"user-authentication/database"
	"user-authentication/helpers"

	"github.com/gin-gonic/gin"
)

func userLogin(ctx *gin.Context) {
	var body helpers.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.String(400, "Invalid request %s", err)
		return
	}
	email := body.Email
	password := body.Password

	db, err := database.GetDatabaseClient()
	if err != nil {
		ctx.String(500, "Database connection error: %s, %s", err, db)
		return
	}
	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email=$1", email).Scan(&storedPassword)
	if err != nil {
		ctx.String(401, "Invalid username or password")
		return
	}
	if storedPassword != password {
		ctx.String(401, "Invalid password for user %s", email)
		return
	}
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
