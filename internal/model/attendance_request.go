package model

type CreateAttendanceRequest struct {
	StartTime string `json:"start_time" validate:"required,is-valid-datetime"`
	EndTime   string `json:"end_time" validate:"required,is-valid-datetime"`
}
