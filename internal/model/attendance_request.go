package model

// CreateAttendanceRequest represents the request body for creating attendance record
// swagger:model CreateAttendanceRequest
type CreateAttendanceRequest struct {
	// Start time of work shift (YYYY-MM-DD HH:MM:SS format)
	// required: true
	// example: "2024-01-15 08:00:00"
	StartTime string `json:"start_time" validate:"required,is-valid-datetime"`

	// End time of work shift (YYYY-MM-DD HH:MM:SS format)
	// required: true
	// example: "2024-01-15 17:00:00"
	EndTime string `json:"end_time" validate:"required,is-valid-datetime"`
}
