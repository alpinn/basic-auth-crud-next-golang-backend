package services

import (
	"fmt"
	"log"
	"time"

	"github.com/alpinn/auth-go/config"
	models "github.com/alpinn/auth-go/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var Rdb *redis.Client

func InitRedis(client *redis.Client) {
	Rdb = client
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetUsers(db *sqlx.DB) ([]models.User, error) {
	rows, err := db.Query("SELECT id, name, email, role, created_at, updated_at FROM public.users")
	if err != nil {
		log.Println("Error DB Query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var id string
		// Scan the row into the variables
		err := rows.Scan(&id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		user.ID, err = uuid.Parse(id)
		if err != nil {
			log.Printf("GetUsers: error parsing UUID: %v", err)
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	return users, nil
}

func RegisterUser(db *sqlx.DB, user models.User) error {
	user.ID = uuid.New()
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (id, name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())",
		user.ID, user.Name, user.Email, hashedPassword, user.Role)
	if err != nil {
		return fmt.Errorf("failed to register user: %v", err)
	}
	return nil
}

func LoginUser(db *sqlx.DB, email, password string) (*models.User, error) {
	var user models.User

	row := db.QueryRow("SELECT id, name, email, password, role, created_at, updated_at FROM public.users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("LoginUser: no user found with this email")
		return nil, fmt.Errorf("no user found with this email")
	}

	if !VerifyPassword(user.Password, password) {
		log.Println("LoginUser: password does not match")
		return nil, fmt.Errorf("password does not match")
	}

	// Store session in Redis
	err = Rdb.Set(config.Ctx, user.ID.String(), user.Email, 30*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set session in Redis: %v", err)
	}

	return &user, nil
}

func UpdateUser(db *sqlx.DB, userID uuid.UUID, name string, email string, password string) error {
	var err error
	if password != "" {
		hashedPassword, err := HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		_, err = db.Exec("UPDATE users SET name = $1, email = $2, password = $3, updated_at = NOW() WHERE id = $4",
			name, email, hashedPassword, userID)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	} else {
		_, err = db.Exec("UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3",
			name, email, userID)
		if err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}
	return nil
}
