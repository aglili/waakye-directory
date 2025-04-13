package provider

import (
	"database/sql"

	"github.com/aglili/waakye-directory/internal/config"
	"github.com/aglili/waakye-directory/internal/handlers"
	"github.com/aglili/waakye-directory/internal/repository/postgres"
)

type Provider struct {
	Cfg           *config.Config
	DB            *sql.DB
	UploadHandler *handlers.UploadHandler
	VendorHandler *handlers.VendorHandler
}

func NewProvider(db *sql.DB, cfg *config.Config) *Provider {
	vendorRepository := postgres.NewVendorRepository(db)
	ratingsRepository := postgres.NewRatingRepository(db)

	vendorHandler := handlers.NewVendorHandler(vendorRepository, ratingsRepository)
	uploadHandler := handlers.NewUploadHandler(cfg.FileUploadPath)

	return &Provider{
		DB:            db,
		VendorHandler: vendorHandler,
		UploadHandler: uploadHandler,
		Cfg:           cfg,
	}
}
