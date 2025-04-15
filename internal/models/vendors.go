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
	ImageURL       string    `json:"image_url" db:"image_url"`
	PhoneNumber    string    `json:"phone_number" db:"phone_number"`
	IsVerified     bool      `json:"is_verified" db:"is_verified"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	Distance       float64   `json:"distance_km,omitempty" db:"-"`
	AverageRating float64   `json:"average_rating" db:"average_rating"`
	AverageHygieneRating float64 `json:"average_hygiene_rating" db:"average_hygiene_rating"`
	AverageValueRating float64 `json:"average_value_rating" db:"average_value_rating"`
	AverageTasteRating float64 `json:"average_taste_rating" db:"average_taste_rating"`
	AverageServiceRating float64 `json:"average_service_rating" db:"average_service_rating"`
	Ratings []VendorRating `json:"ratings" db:"-"`
}




type VendorRating struct {
    ID            uuid.UUID
    HygieneRating int
    ValueRating   int
    ServiceRating int
    TasteRating   int
    Comment       string
    CreatedAt     time.Time
}