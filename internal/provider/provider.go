package provider

import (
	"database/sql"

	"github.com/aglili/waakye-directory/internal/handlers"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
)

type Provider struct {
	DB            *sql.DB
	VendorHandler *handlers.VendorHandler
}

func NewProvider(db *sql.DB) *Provider {
	vendorRepository := postgres.NewVendorRepository(db)
	ratingsRepository := postgres.NewRatingRepository(db)

	vendorHandler := handlers.NewVendorHandler(vendorRepository,ratingsRepository)

	return &Provider{
		DB:            db,
		VendorHandler: vendorHandler,
	}
}
