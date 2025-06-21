package model

// SignInRequest represents the request body for user authentication
// swagger:model SignInRequest
type SignInRequest struct {
	// Username for authentication (4-100 characters)
	// required: true
	// min: 4
	// max: 100
	// example: "john.doe"
	Username string `json:"username" validate:"required,min=4,max=100"`

	// Password for authentication (must be strong password)
	// required: true
	// example: "StrongP@ssw0rd123"
	Password string `json:"password" validate:"required,is-strong-password"`
}

// VerifyAccountRequest represents the request body for account verification
// swagger:model VerifyAccountRequest
type VerifyAccountRequest struct {
	// Verification token
	// required: true
	// max: 100
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	Token string `validate:"required,max=100"`
}
