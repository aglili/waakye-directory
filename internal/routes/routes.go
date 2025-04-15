package routes

import (
	"net/http"
	"path/filepath"

	"github.com/aglili/waakye-directory/internal/provider"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	swaggerfiles "github.com/swaggo/files"
	_ "github.com/aglili/waakye-directory/docs"
)




func SetupRoutes(provider *provider.Provider) http.Handler {
	router := gin.Default()

	uploadsDir := filepath.Join(provider.Cfg.FileUploadPath)
	router.StaticFS("/uploads", http.Dir(uploadsDir))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := router.Group("/api/v1")

	v1.POST("/vendors", provider.VendorHandler.CreateVendor)
	v1.GET("/vendors", provider.VendorHandler.ListVendorsWithPagination)
	v1.GET("/vendors/:id", provider.VendorHandler.GetVendorByID)
	v1.GET("/vendors/nearby", provider.VendorHandler.GetNearbyVendors)
	v1.GET("/vendors/verified", provider.VendorHandler.GetVerifiedVendors)
	v1.GET("/vendors/top_rated", provider.VendorHandler.GetTopRatedVendors)
	v1.POST("/vendors/:id/rate", provider.VendorHandler.RateVendor)
	v1.GET("/vendors/:id/ratings", provider.VendorHandler.GetVendorRatings)

	v1.POST("/uploads", provider.UploadHandler.UploadFile)
	return router
}
