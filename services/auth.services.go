package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alpinn/auth-go/config"
	models "github.com/alpinn/auth-go/models"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

var RedisClient *redis.Client

func InitRedis(client *redis.Client) {
	RedisClient = client
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(hashedPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	return err == nil
}

func RegisterUser(db *sql.DB, user models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (:1, :2, :3, :4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", user.Name, user.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}
	return nil
}

func LoginUser(db *sql.DB, email, password string) (*models.User, error) {
	var user models.User
	row := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = :1", email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	if !VerifyPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Store session in Redis
	err = RedisClient.Set(config.Ctx, fmt.Sprintf("session:%d", user.ID), user.Email, 30*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set session in Redis: %v", err)
	}

	return &user, nil
}
