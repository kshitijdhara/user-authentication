package helpers

import (
	"user-authentication/database"
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

func CreateUser(email, firstName, lastName, password, userType string) error {
	db, err := database.GetDatabaseClient()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, first_name, last_name, email, user_type, image, password, version, created_at, updated_at, permission, is_verified) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, 1, NOW(), NOW(), 'user', false)",
		firstName, lastName, email, userType, "{}", password)
	return err
}

func CreateJWTToken(userId string) (string, error) {
	return "mock-jwt-token-for-" + userId, nil
}
