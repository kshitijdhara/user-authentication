package helpers

import (
	"fmt"
	"os"
	"time"
	"user-authentication/database"

	"github.com/golang-jwt/jwt/v5"
)

func NewUserCheck(email string) (string, bool, error) {
	db, err := database.GetDatabaseClient()
	if err != nil {
		return "", false, err
	}
	var userId string
	err = db.QueryRow("SELECT id FROM users WHERE email=$1", email).Scan(&userId)
	if err != nil {
		return "", false, nil
	}
	return userId, true, nil
}

func CreateUser(email, firstName, lastName, password, userType string) (string, error) {
	db, err := database.GetDatabaseClient()
	if err != nil {
		return "", err
	}
	var userId string
	err = db.QueryRow("INSERT INTO users (id, first_name, last_name, email, user_type, image, password, version, created_at, updated_at, permission, is_verified) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, 1, NOW(), NOW(), 'user', false) RETURNING id",
		firstName, lastName, email, userType, "{}", password).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func CreateJWTToken(userId string) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", fmt.Errorf("JWT_KEY environment variable not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWTToken(tokenString string) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", fmt.Errorf("JWT_KEY environment variable not set")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["user_id"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token claims")
		}
		return userId, nil
	}
	return "", fmt.Errorf("invalid token")
}
