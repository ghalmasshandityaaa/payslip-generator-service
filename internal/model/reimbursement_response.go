package model

// CreateReimbursementRequest represents the request body for creating reimbursement record
// swagger:model CreateReimbursementRequest
type CreateReimbursementRequest struct {
	// Reimbursement amount in currency units
	// required: true
	// minimum: 1
	// example: 150000
	Amount int `json:"amount" validate:"required,min=1"`

	// Description of the reimbursement expense
	// required: true
	// min: 5
	// max: 255
	// example: "Transportation expenses for client meeting"
	Description string `json:"description" validate:"required,min=5,max=255"`
}
