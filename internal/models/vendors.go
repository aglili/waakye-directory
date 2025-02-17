package models

import (
	"github.com/google/uuid"
	"time"
)

type Location struct {
	ID            uuid.UUID `json:"id" db:"id"`
	StreetAddress string    `json:"street_address" db:"street_address"`
	City          string    `json:"city" db:"city"`
	Region        string    `json:"region" db:"region"`
	Latitude      float64   `json:"latitude" db:"latitude"`
	Longitude     float64   `json:"longitude" db:"longitude"`
	Landmark      string    `json:"landmark" db:"landmark"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type WaakyeVendor struct {
    ID             uuid.UUID `json:"id" db:"id"`
    Name           string    `json:"name" db:"name"`
    LocationID     uuid.UUID `json:"location_id" db:"location_id"`
    Location       Location  `json:"location" db:"-"`
    Description    string    `json:"description" db:"description"`
    OperatingHours string    `json:"operating_hours" db:"operating_hours"`
    PhoneNumber    string    `json:"phone_number" db:"phone_number"`
    IsVerified     bool      `json:"is_verified" db:"is_verified"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
    Distance       float64   `json:"distance_km,omitempty" db:"-"`
}
