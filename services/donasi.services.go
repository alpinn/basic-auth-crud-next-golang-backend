package services

import (
	"fmt"
	"log"
	"strings"

	models "github.com/alpinn/auth-go/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func PostDonasi(db *sqlx.DB, donasi models.Donasi) error {
	donasi.ID = uuid.New()
	uuidHexID := strings.ReplaceAll(donasi.ID.String(), "-", "")
	uuidHexUserID := strings.ReplaceAll(donasi.UserID.String(), "-", "")

	_, err := db.Exec("INSERT INTO donasi (id, user_id, name, nominal, pesan, url, created_at, updated_at) VALUES (:1, :2, :3, :4, :5, :6, SYSDATE, SYSDATE)",
		uuidHexID, uuidHexUserID, donasi.Name, donasi.Nominal, donasi.Pesan, donasi.Url)
	if err != nil {
		log.Println("Error inserting donation:", err)
		return fmt.Errorf("failed to create donasi: %v", err)
	}
	return nil
}

func GetDonasi(db *sqlx.DB) ([]models.Donasi, error) {
	rows, err := db.Query("SELECT d.id, d.user_id, d.nominal, d.pesan, d.url, d.created_at, d.updated_at, u.name FROM donasi d JOIN users u ON d.user_id = u.id")
	if err != nil {
		log.Println("Error querying donations:", err)
		return nil, fmt.Errorf("failed to get donations: %v", err)
	}
	defer rows.Close()

	var donasis []models.Donasi
	for rows.Next() {
		var donasi models.Donasi
		if err := rows.Scan(&donasi.ID, &donasi.UserID, &donasi.Nominal, &donasi.Pesan, &donasi.Url, &donasi.CreatedAt, &donasi.UpdatedAt, &donasi.Name); err != nil {
			log.Println("Error scanning donation row:", err)
			return nil, err
		}
		donasis = append(donasis, donasi)
	}

	return donasis, nil
}
