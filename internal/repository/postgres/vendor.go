package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aglili/waakye-directory/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type VendorRepository interface {
	CreateVendor(ctx context.Context, vendor *models.WaakyeVendor) error
	ListVendorsWithPagination(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error)
	CountVendors(ctx context.Context) (int64, error)
	GetVendorByID(ctx context.Context, id uuid.UUID) (*models.WaakyeVendor, error)
	GetNearbyVendors(ctx context.Context, latitude, longitude, radius float64) ([]models.WaakyeVendor, error)
	GetVerifiedVendors(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error)
	CountVerifiedVendors(ctx context.Context) (int64, error)
	GetTopRatedVendors(ctx context.Context) ([]models.WaakyeVendor, error)
}

type vendorRepository struct {
	db *sql.DB
}

func NewVendorRepository(db *sql.DB) VendorRepository {
	return &vendorRepository{
		db: db,
	}
}

func (r *vendorRepository) CreateVendor(ctx context.Context, vendor *models.WaakyeVendor) error {
	query := `
		WITH location_insert AS (
			INSERT INTO locations (street_address, city, region, latitude, longitude, landmark)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		)
		INSERT INTO waakye_vendors (name, location_id, description, operating_hours,image_url, phone_number, is_verified)
		SELECT $7, id, $8, $9, $10, $11,$12
		FROM location_insert
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		vendor.Location.StreetAddress,
		vendor.Location.City,
		vendor.Location.Region,
		vendor.Location.Latitude,
		vendor.Location.Longitude,
		vendor.Location.Landmark,
		vendor.Name,
		vendor.Description,
		vendor.OperatingHours,
		vendor.ImageURL,
		vendor.PhoneNumber,
		vendor.IsVerified,
	).Scan(&vendor.ID, &vendor.CreatedAt, &vendor.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("Failed to create vendor: no rows returned")
			return errors.New("failed to create vendor: no rows returned")
		}
		log.Error().Err(err).Msg("Failed to create vendor")
		return err
	}

	return nil
}

func (r *vendorRepository) ListVendorsWithPagination(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error) {
	query := `
		SELECT wv.id, wv.name, wv.description, wv.operating_hours,wv.image_url, wv.phone_number, wv.is_verified, wv.created_at, wv.updated_at,
			l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
		FROM waakye_vendors wv
		INNER JOIN locations l ON wv.location_id = l.id
		ORDER BY wv.created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list vendors")
		return nil, err
	}

	defer rows.Close()

	var vendors []models.WaakyeVendor
	for rows.Next() {
		var vendor models.WaakyeVendor
		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.Description,
			&vendor.OperatingHours,
			&vendor.ImageURL,
			&vendor.PhoneNumber,
			&vendor.IsVerified,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
			&vendor.Location.StreetAddress,
			&vendor.Location.City,
			&vendor.Location.Region,
			&vendor.Location.Latitude,
			&vendor.Location.Longitude,
			&vendor.Location.Landmark,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan vendor")
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	if len(vendors) == 0 {
		return []models.WaakyeVendor{}, nil
	}

	return vendors, nil

}

func (r *vendorRepository) CountVendors(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM waakye_vendors`

	var totalItems int64
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count vendors")
		return 0, err
	}

	return totalItems, nil

}

func (r *vendorRepository) GetVendorByID(ctx context.Context, id uuid.UUID) (*models.WaakyeVendor, error) {
    // First, get the vendor details
    vendorQuery := `
        SELECT wv.id, wv.name, wv.description, wv.operating_hours, wv.image_url, wv.phone_number, 
               wv.is_verified, wv.created_at, wv.updated_at,
               l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark,
               COALESCE(AVG((vr.hygiene_rating + vr.value_rating + vr.taste_rating + vr.service_rating) / 4), 0) as avg_rating,
               COALESCE(AVG(vr.hygiene_rating), 0) as avg_hygiene_rating,
               COALESCE(AVG(vr.value_rating), 0) as avg_value_rating,
               COALESCE(AVG(vr.taste_rating), 0) as avg_taste_rating,
               COALESCE(AVG(vr.service_rating), 0) as avg_service_rating
        FROM waakye_vendors wv
        INNER JOIN locations l ON wv.location_id = l.id
        LEFT JOIN vendor_ratings vr ON wv.id = vr.vendor_id
        WHERE wv.id = $1
        GROUP BY wv.id, l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
    `

    var vendor models.WaakyeVendor
    err := r.db.QueryRowContext(ctx, vendorQuery, id).Scan(
        &vendor.ID,
        &vendor.Name,
        &vendor.Description,
        &vendor.OperatingHours,
        &vendor.ImageURL,
        &vendor.PhoneNumber,
        &vendor.IsVerified,
        &vendor.CreatedAt,
        &vendor.UpdatedAt,
        &vendor.Location.StreetAddress,
        &vendor.Location.City,
        &vendor.Location.Region,
        &vendor.Location.Latitude,
        &vendor.Location.Longitude,
        &vendor.Location.Landmark,
        &vendor.AverageRating,
        &vendor.AverageHygieneRating,
        &vendor.AverageValueRating,
        &vendor.AverageTasteRating,
        &vendor.AverageServiceRating,
    )
    if err != nil {
        log.Error().Err(err).Msg("Failed to get vendor by ID")
        return nil, err
    }

    // Now get comments for this vendor
    commentsQuery := `
        SELECT id, hygiene_rating, value_rating, service_rating, taste_rating, comment, created_at
        FROM vendor_ratings
        WHERE vendor_id = $1
        ORDER BY created_at DESC
    `

    commentsRows, err := r.db.QueryContext(ctx, commentsQuery, id)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get comments for vendor")
        return nil, err
    }
    defer commentsRows.Close()

    var comments []models.VendorRating
    for commentsRows.Next() {
        var rating models.VendorRating
        err := commentsRows.Scan(
            &rating.ID,
            &rating.HygieneRating,
            &rating.ValueRating,
            &rating.ServiceRating,
            &rating.TasteRating,
            &rating.Comment,
            &rating.CreatedAt,
        )
        if err != nil {
            log.Error().Err(err).Msg("Failed to scan comment")
            return nil, err
        }

        comments = append(comments, rating)
    }

    // Add comments to the vendor
    vendor.Ratings = comments

    return &vendor, nil
}

func (r *vendorRepository) GetNearbyVendors(ctx context.Context, latitude, longitude, radiusKm float64) ([]models.WaakyeVendor, error) {
    // Convert radius from kilometers to meters
    radiusMeters := radiusKm * 1000.0

    query := `
        SELECT 
            wv.id, 
            wv.name, 
            wv.description, 
            wv.operating_hours, 
            wv.image_url,
            wv.phone_number, 
            wv.is_verified, 
            wv.created_at, 
            wv.updated_at,
            l.street_address, 
            l.city, 
            l.region, 
            l.latitude, 
            l.longitude, 
            l.landmark,
            earth_distance(ll_to_earth($1, $2), ll_to_earth(l.latitude, l.longitude)) as distance,
            COALESCE(AVG((vr.hygiene_rating + vr.value_rating + vr.taste_rating + vr.service_rating) / 4), 0) as avg_rating
        FROM waakye_vendors wv
        INNER JOIN locations l ON wv.location_id = l.id
        LEFT JOIN vendor_ratings vr ON wv.id = vr.vendor_id
        WHERE earth_distance(ll_to_earth($1, $2), ll_to_earth(l.latitude, l.longitude)) <= $3
        GROUP BY wv.id, l.id
        ORDER BY distance ASC
    `

    rows, err := r.db.QueryContext(ctx, query, latitude, longitude, radiusMeters)
    if err != nil {
        log.Error().Err(err).
            Float64("latitude", latitude).
            Float64("longitude", longitude).
            Float64("radius_km", radiusKm).
            Msg("Failed to get nearby vendors")
        return nil, err
    }
    defer rows.Close()

    var vendors []models.WaakyeVendor
    for rows.Next() {
        var vendor models.WaakyeVendor
        var distance float64
        var avgRating float64
        err := rows.Scan(
            &vendor.ID,
            &vendor.Name,
            &vendor.Description,
            &vendor.OperatingHours,
            &vendor.ImageURL,
            &vendor.PhoneNumber,
            &vendor.IsVerified,
            &vendor.CreatedAt,
            &vendor.UpdatedAt,
            &vendor.Location.StreetAddress,
            &vendor.Location.City,
            &vendor.Location.Region,
            &vendor.Location.Latitude,
            &vendor.Location.Longitude,
            &vendor.Location.Landmark,
            &distance,
            &avgRating,
        )
        if err != nil {
            log.Error().Err(err).Msg("Failed to scan vendor")
            return nil, err
        }

        // Convert distance from meters to kilometers and add to vendor
        vendor.Distance = distance / 1000.0
        
        // Set the average rating
        vendor.AverageRating = avgRating

        vendors = append(vendors, vendor)
    }

    if len(vendors) == 0 {
        return []models.WaakyeVendor{}, nil
    }

    return vendors, nil
}
func (r *vendorRepository) GetVerifiedVendors(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error) {
	query := `
		SELECT wv.id, wv.name, wv.description, wv.operating_hours,wv.image_url ,wv.phone_number, wv.is_verified, wv.created_at, wv.updated_at,
			l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
		FROM waakye_vendors wv
		INNER JOIN locations l ON wv.location_id = l.id
		WHERE wv.is_verified = true
		ORDER BY wv.created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list vendors")
		return nil, err
	}

	defer rows.Close()

	var vendors []models.WaakyeVendor

	for rows.Next() {
		var vendor models.WaakyeVendor
		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.Description,
			&vendor.OperatingHours,
			&vendor.ImageURL,
			&vendor.PhoneNumber,
			&vendor.IsVerified,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
			&vendor.Location.StreetAddress,
			&vendor.Location.City,
			&vendor.Location.Region,
			&vendor.Location.Latitude,
			&vendor.Location.Longitude,
			&vendor.Location.Landmark,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan vendor")
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	if len(vendors) == 0 {
		return []models.WaakyeVendor{}, nil
	}

	return vendors, nil

}

func (r *vendorRepository) CountVerifiedVendors(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM waakye_vendors WHERE is_verified = true`

	var totalItems int64
	err := r.db.QueryRowContext(ctx, query).Scan(&totalItems)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count vendors")
		return 0, err
	}

	return totalItems, nil
}

func (r *vendorRepository) GetTopRatedVendors(ctx context.Context) ([]models.WaakyeVendor, error) {
	query := `
		SELECT wv.id, wv.name, wv.description, wv.operating_hours,wv.image_url ,wv.phone_number, wv.is_verified, wv.created_at, wv.updated_at,
			l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
		FROM waakye_vendors wv
		INNER JOIN locations l ON wv.location_id = l.id
		INNER JOIN vendor_ratings vr ON wv.id = vr.vendor_id
		GROUP BY wv.id, l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
		ORDER BY AVG((vr.hygiene_rating + vr.value_rating + vr.taste_rating + vr.service_rating) / 4) DESC
		LIMIT 5
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list vendors")
		return nil, err
	}

	defer rows.Close()

	var vendors []models.WaakyeVendor

	for rows.Next() {
		var vendor models.WaakyeVendor
		err := rows.Scan(
			&vendor.ID,
			&vendor.Name,
			&vendor.Description,
			&vendor.OperatingHours,
			&vendor.ImageURL,
			&vendor.PhoneNumber,
			&vendor.IsVerified,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
			&vendor.Location.StreetAddress,
			&vendor.Location.City,
			&vendor.Location.Region,
			&vendor.Location.Latitude,
			&vendor.Location.Longitude,
			&vendor.Location.Landmark,
		)

		if err != nil {
			log.Error().Err(err).Msg("Failed to scan vendor")
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	if len(vendors) == 0 {
		return []models.WaakyeVendor{}, nil
	}

	return vendors, nil
}
