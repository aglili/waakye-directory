package handlers





type UploadResponse struct {
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	FileType string `json:"file_type"`
}



type BadRequestResponse struct {
	Error  string `json:"error"`
	Details string `json:"details"`
}


type InternalServerErrorResponse struct {
	Error  string `json:"error"`
	Details string `json:"details"`
}


type CreatedResponse struct{
	Data    map[string]interface{}   `json:"data"`
	Message string      `json:"message"`

}

type LocationSchema struct {
	StreetAddress string  `json:"street_address" binding:"required"`
	City          string  `json:"city" binding:"required"`
	Region        string  `json:"region" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	Landmark      string  `json:"landmark" binding:"required"`
}

type CreateWaakyeVendorSchema struct {
	Name           string        `json:"name" binding:"required"`
	Location       LocationSchema `json:"location" binding:"required"`
	Description    string        `json:"description" binding:"required"`
	OperatingHours string        `json:"operating_hours" binding:"required"`
	ImageURL       string        `json:"image_url" binding:"required"`
	PhoneNumber    string        `json:"phone_number" binding:"required"`
}



type PaginatedResponse struct{
	Data       []map[string]interface{} `json:"data"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalItems int64                  `json:"total_items"`
	TotalPages int                     `json:"total_pages"`
	Message    string                  `json:"message"`
}


type NotFoundResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}



type RateVendorRequest struct {
	HygieneRating int    `json:"hygeine_rating" binding:"required"`
	ValueRating   int    `json:"value_rating" validate:"required,gte=1,lte=5" db:"value_rating"`
	TasteRating   int    `json:"taste_rating" validate:"required,gte=1,lte=5" db:"taste_rating"`
	ServiceRating int    `json:"service_rating" validate:"required,gte=1,lte=5" db:"service_rating"`
	Comment       string `json:"comment" db:"comment"`
}