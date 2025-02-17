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

	vendorHandler := handlers.NewVendorHandler(vendorRepository)

	return &Provider{
		DB:            db,
		VendorHandler: vendorHandler,
	}
}
