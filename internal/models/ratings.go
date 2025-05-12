package models

import "time"

type RateVendorRequest struct {
	HygeineRating int    `json:"hygeine_rating" validate:"required,gte=1,lte=5" db:"hygeine_rating"`
	ValueRating   int    `json:"value_rating" validate:"required,gte=1,lte=5" db:"value_rating"`
	TasteRating   int    `json:"taste_rating" validate:"required,gte=1,lte=5" db:"taste_rating"`
	ServiceRating int    `json:"service_rating" validate:"required,gte=1,lte=5" db:"service_rating"`
	Comment       string `json:"comment" db:"comment"`
}

type VendorRatings struct {
	HygieneRating float64         `json:"hygiene_rating" db:"hygiene_rating"`
	ValueRating   float64         `json:"value_rating" db:"value_rating"`
	TasteRating   float64         `json:"taste_rating" db:"taste_rating"`
	ServiceRating float64         `json:"service_rating" db:"service_rating"`
	OverallRating float64         `json:"overall_rating" db:"overall_rating"`
	TotalRatings  int             `json:"total_ratings" db:"total_ratings"`
	Comments      []VendorComment `json:"comments"`
}

// New struct to hold individual comments
type VendorComment struct {
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
