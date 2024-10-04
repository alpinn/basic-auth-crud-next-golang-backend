package models

import (
	"time"

	"github.com/google/uuid"
)

type Donasi struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Nominal   int       `json:"nominal"`
	Pesan     string    `json:"pesan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
