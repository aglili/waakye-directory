package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aglili/waakye-directory/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



func setupMockDB(t *testing.T) (*sql.DB,sqlmock.Sqlmock,VendorRepository){
	db, mock, err := sqlmock.New()
	require.NoError(t, err)


	repo := NewVendorRepository(db)
	return db, mock, repo
}


func TestCreateVendor(t *testing.T) {
    db, mock, repo := setupMockDB(t)
    defer db.Close()
    
    ctx := context.Background()
    testID := uuid.New()
    testVendorVerfied := true
    now := time.Now()
    
    vendor := &models.WaakyeVendor{
        Name: "Test Vendor",
        Location: models.Location{
            StreetAddress: "Test Street",
            City: "Test City",
            Region: "Test Region",
            Latitude: 5.0,
            Longitude: 5.0,
            Landmark: "Test Landmark",
        },
        Description: "Test Description",
        OperatingHours: "Test Operating Hours",
        PhoneNumber: "Test Phone Number",
    }
    
    mock.ExpectQuery(`WITH location_insert AS`).
        WithArgs(
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
        ).
        WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
            AddRow(testID, now, now))
    
    err := repo.CreateVendor(ctx, vendor)
    require.NoError(t, err)
    require.Equal(t, testID, vendor.ID)
    require.Equal(t, now, vendor.CreatedAt)
    require.Equal(t, now, vendor.UpdatedAt)
    require.NoError(t, mock.ExpectationsWereMet())
    assert.NotEqual(t,testVendorVerfied,vendor.IsVerified)
}



func TestListVendorsWithPagination(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	page, pageSize := 1, 10
	
	vendorID1 := uuid.New()
	vendorID2 := uuid.New()
	now := time.Now()

	// Expected rows to be returned by the query
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "operating_hours", "phone_number", "is_verified", "created_at", "updated_at",
		"street_address", "city", "region", "latitude", "longitude", "landmark",
	}).
		AddRow(
			vendorID1, "Vendor One", "Description 1", "8am - 5pm", "+233201234567", true, now, now,
			"123 Street", "Accra", "Greater Accra", 5.6037, -0.1870, "Landmark 1",
		).
		AddRow(
			vendorID2, "Vendor Two", "Description 2", "9am - 6pm", "+233209876543", false, now, now,
			"456 Street", "Kumasi", "Ashanti", 6.7000, -1.6167, "Landmark 2",
		)


	mock.ExpectQuery(`SELECT wv.id, wv.name, wv.description`).
		WithArgs(pageSize, 0). // LIMIT 10 OFFSET 0
		WillReturnRows(rows)


	vendors, err := repo.ListVendorsWithPagination(ctx, page, pageSize)
	assert.NoError(t, err)
	assert.Len(t, vendors, 2)
	
	assert.Equal(t, "Vendor One", vendors[0].Name)
	assert.Equal(t, "Vendor Two", vendors[1].Name)
	assert.Equal(t, "+233201234567", vendors[0].PhoneNumber)
	assert.Equal(t, "Kumasi", vendors[1].Location.City)

	assert.NoError(t, mock.ExpectationsWereMet())
}



func TestCountVendors(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM waakye_vendors`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))


	count, err := repo.CountVendors(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)

	assert.NoError(t, mock.ExpectationsWereMet())
}



func TestGetVendorByID(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	testID := uuid.New()
	now := time.Now()
	

	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "operating_hours", "phone_number", "is_verified", "created_at", "updated_at",
		"street_address", "city", "region", "latitude", "longitude", "landmark",
	}).
		AddRow(
			testID, "Test Vendor", "Description", "8am - 5pm", "+233201234567", true, now, now,
			"123 Street", "Accra", "Greater Accra", 5.6037, -0.1870, "Near market",
		)

	mock.ExpectQuery(`SELECT wv.id, wv.name, wv.description`).
		WithArgs(testID).
		WillReturnRows(rows)

	vendor, err := repo.GetVendorByID(ctx, testID)
	assert.NoError(t, err)
	assert.NotNil(t, vendor)
	assert.Equal(t, testID, vendor.ID)
	assert.Equal(t, "Test Vendor", vendor.Name)
	assert.Equal(t, "Accra", vendor.Location.City)

	assert.NoError(t, mock.ExpectationsWereMet())
}




func TestGetNearbyVendors(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	ctx := context.Background()
	latitude, longitude := 5.6037, -0.1870
	radiusKm := 5.0
	vendorID := uuid.New()
	now := time.Now()
	
    
	rows := sqlmock.NewRows([]string{
		"id", "name", "description", "operating_hours", "phone_number", "is_verified", "created_at", "updated_at",
		"street_address", "city", "region", "latitude", "longitude", "landmark", "distance",
	}).
		AddRow(
			vendorID, "Nearby Vendor", "Close by", "9am - 6pm", "+233207654321", true, now, now,
			"10 Close Street", "Accra", "Greater Accra", 5.6057, -0.1890, "Corner shop", 500.0,
		)

	mock.ExpectQuery(`SELECT(.*)earth_distance(.*)ll_to_earth`).
		WithArgs(latitude, longitude, radiusKm*1000.0).
		WillReturnRows(rows)

	// Test the GetNearbyVendors method
	vendors, err := repo.GetNearbyVendors(ctx, latitude, longitude, radiusKm)
	assert.NoError(t, err)
	assert.Len(t, vendors, 1)
	assert.Equal(t, "Nearby Vendor", vendors[0].Name)
	assert.Equal(t, 0.5, vendors[0].Distance) // 500m converted to 0.5km

	assert.NoError(t, mock.ExpectationsWereMet())
}