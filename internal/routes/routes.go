package routes

import (
	"net/http"

	"github.com/aglili/waakye-directory/internal/provider"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(provider *provider.Provider) http.Handler {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/vendors", provider.VendorHandler.CreateVendor)
	v1.GET("/vendors", provider.VendorHandler.ListVendorsWithPagination)
	v1.GET("/vendors/:id", provider.VendorHandler.GetVendorByID)
	v1.GET("/vendors/nearby", provider.VendorHandler.GetNearbyVendors)

	return router
}
