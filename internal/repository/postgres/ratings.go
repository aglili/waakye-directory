package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aglili/waakye-directory/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)




type RatingsRepository interface {
	RateVendor(ctx context.Context, request *models.RateVendorRequest) error
	GetVendorRatings(ctx context.Context, vendorID uuid.UUID) (*models.VendorRatings, error)
}

type ratingsRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) RatingsRepository {
	return &ratingsRepository{
		db: db,
	}
}




func (r *ratingsRepository) RateVendor(ctx context.Context, request *models.RateVendorRequest) error {
	query := `
		INSERT INTO vendor_ratings (vendor_id, hygiene_rating, value_rating, taste_rating, service_rating, comment)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		request.VendorID,
		request.HygeineRating,
		request.ValueRating,
		request.TasteRating,
		request.ServiceRating,
		request.Comment,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to rate vendor")
		return errors.New("failed to rate vendor")
	}

	return nil
}


func (r *ratingsRepository) GetVendorRatings(ctx context.Context, vendorID uuid.UUID) (*models.VendorRatings, error) {
	query := `
		SELECT
			COALESCE(AVG(hygiene_rating), 0) AS hygiene_rating,
			COALESCE(AVG(value_rating), 0) AS value_rating,
			COALESCE(AVG(taste_rating), 0) AS taste_rating,
			COALESCE(AVG(service_rating), 0) AS service_rating,
			COALESCE(AVG((hygiene_rating + value_rating + taste_rating + service_rating) / 4), 0) AS overall_rating,
			COUNT(*) AS total_ratings
		FROM vendor_ratings
		WHERE vendor_id = $1
	`

	var ratings models.VendorRatings
	err := r.db.QueryRowContext(ctx, query, vendorID).Scan(
		&ratings.HygieneRating,
		&ratings.ValueRating,
		&ratings.TasteRating,
		&ratings.ServiceRating,
		&ratings.OverallRating,
		&ratings.TotalRatings,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.VendorRatings{}, nil
		}
		log.Error().Err(err).Msg("Failed to get vendor ratings")
		return nil, fmt.Errorf("failed to get vendor ratings: %w", err)
	}

	return &ratings, nil
}

