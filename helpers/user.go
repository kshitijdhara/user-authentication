package helpers

import (
	"crypto/sha256"
	"encoding/hex"
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

func CreateTokenPair(userId, role string) (string, string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", "", fmt.Errorf("JWT_KEY environment variable not set")
	}

	// Access Token (15 mins)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7 days)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(key))
	if err != nil {
		return "", "", err
	}

	// Store Refresh Token in DB
	db, err := database.GetDatabaseClient()
	if err != nil {
		return "", "", err
	}

	// Hash the refresh token
	hash := sha256.Sum256([]byte(signedRefreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	_, err = db.Exec("INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		userId, tokenHash, time.Now().Add(time.Hour*24*7))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func ValidateJWTToken(tokenString string) (string, string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", "", fmt.Errorf("JWT_KEY environment variable not set")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Ensure it's not a refresh token being used as access token
		if claims["type"] == "refresh" {
			return "", "", fmt.Errorf("invalid token type")
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid token claims: user_id")
		}
		role, ok := claims["role"].(string)
		if !ok {
			role = "user"
		}
		return userId, role, nil
	}
	return "", "", fmt.Errorf("invalid token")
}

func ValidateRefreshToken(tokenString string) (string, string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", "", fmt.Errorf("JWT_KEY environment variable not set")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "refresh" {
			return "", "", fmt.Errorf("invalid token type")
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			return "", "", fmt.Errorf("invalid token claims: user_id")
		}
		role, ok := claims["role"].(string)
		if !ok {
			role = "user"
		}

		// Check DB
		db, err := database.GetDatabaseClient()
		if err != nil {
			return "", "", err
		}

		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])

		var revoked bool
		err = db.QueryRow("SELECT revoked FROM refresh_tokens WHERE token_hash=$1 AND expires_at > NOW()", tokenHash).Scan(&revoked)
		if err != nil {
			return "", "", fmt.Errorf("invalid or expired refresh token")
		}
		if revoked {
			return "", "", fmt.Errorf("refresh token revoked")
		}

		return userId, role, nil
	}
	return "", "", fmt.Errorf("invalid token")
}

func RevokeRefreshToken(tokenString string) error {
	db, err := database.GetDatabaseClient()
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	_, err = db.Exec("UPDATE refresh_tokens SET revoked=true WHERE token_hash=$1", tokenHash)
	return err
}
