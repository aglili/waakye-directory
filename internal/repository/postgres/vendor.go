package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aglili/waakye-directory/internal/models"
)

type VendorRepository interface {
	CreateVendor(ctx context.Context, vendor *models.WaakyeVendor) error
	ListVendorsWithPagination(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error)
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
		INSERT INTO waakye_vendors (name, location_id, description, operating_hours, phone_number, is_verified)
		SELECT $7, id, $8, $9, $10, $11
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
		vendor.PhoneNumber,
		vendor.IsVerified,
	).Scan(&vendor.ID, &vendor.CreatedAt, &vendor.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("failed to create vendor: no rows returned")
		}
		return err
	}

	return nil
}

func (r *vendorRepository) ListVendorsWithPagination(ctx context.Context, page, pageSize int) ([]models.WaakyeVendor, error) {
	query := `
		SELECT wv.id, wv.name, wv.description, wv.operating_hours, wv.phone_number, wv.is_verified, wv.created_at, wv.updated_at,
			l.street_address, l.city, l.region, l.latitude, l.longitude, l.landmark
		FROM waakye_vendors wv
		INNER JOIN locations l ON wv.location_id = l.id
		ORDER BY wv.created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
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
			return nil, err
		}

		vendors = append(vendors, vendor)
	}

	if len(vendors) == 0 {
		return []models.WaakyeVendor{}, nil
	}

	return vendors, nil

}
