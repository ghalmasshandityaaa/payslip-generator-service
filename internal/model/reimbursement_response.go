package model

type CreateReimbursementRequest struct {
	Amount      int    `json:"amount" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=5,max=255"`
}
