package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alpinn/auth-go/config"
	models "github.com/alpinn/auth-go/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var Rdb *redis.Client

func InitRedis(client *redis.Client) {
	Rdb = client
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(hashedPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	return err == nil
}

func GetUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT ID, NAME, EMAIL, CREATED_AT, UPDATED_AT FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id string
		err := rows.Scan(&id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		user.ID, err = uuid.Parse(id)
		if err != nil {
			log.Printf("GetUsers: error parsing UUID: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func RegisterUser(db *sql.DB, user models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES (:1, :2, :3, :4, SYSTIMESTAMP, SYSTIMESTAMP)", user.Name, user.Email, hashedPassword)
	log.Printf("Inserting user: name=%s, email=%s, password=%s", user.Name, user.Email, hashedPassword)
	if err != nil {
		log.Printf("Database error: %v", err)
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
	err = Rdb.Set(config.Ctx, fmt.Sprintf("session:%d", user.ID), user.Email, 30*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set session in Redis: %v", err)
	}

	return &user, nil
}
