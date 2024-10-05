package services

import (
	"database/sql"
	"fmt"

	models "github.com/alpinn/auth-go/models"
	"github.com/google/uuid"
)

func PostDonasi(db *sql.DB, donasi models.Donasi) error {
	donasi.ID = uuid.New() // Generate a new UUID for the product
	_, err := db.Exec("INSERT INTO donations (id, user_id, nominal, pesan, created_at, updated_at) VALUES (:1, :2, :3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		donasi.ID.String(), donasi.UserID.String(), donasi.Nominal, donasi.Pesan)
	if err != nil {
		return fmt.Errorf("failed to create donasi: %v", err)
	}
	return nil
}

func GetDonasi(db *sql.DB) ([]models.Donasi, error) {
	rows, err := db.Query("SELECT id, user_id, nominal, pesan, created_at, updated_at FROM donations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var donasis []models.Donasi
	for rows.Next() {
		var donasi models.Donasi
		if err := rows.Scan(&donasi.ID, &donasi.UserID, &donasi.Nominal, &donasi.Pesan, &donasi.CreatedAt, &donasi.UpdatedAt); err != nil {
			return nil, err
		}
		donasis = append(donasis, donasi)
	}
	return donasis, nil
}
