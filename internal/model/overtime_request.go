package model

// CreateOvertimeRequest represents the request body for creating overtime record
// swagger:model CreateOvertimeRequest
type CreateOvertimeRequest struct {
	// Date of overtime work (YYYY-MM-DD format)
	// required: true
	// example: "2024-01-15"
	Date string `json:"date" validate:"required,is-valid-date"`

	// Total hours of overtime (1-3 hours maximum)
	// required: true
	// minimum: 1
	// maximum: 3
	// example: 2
	TotalHours int `json:"total_hours" validate:"required,min=1,max=3"`
}
