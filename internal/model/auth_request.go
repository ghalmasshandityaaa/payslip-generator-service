package model

type SignInRequest struct {
	Username string `json:"username" validate:"required,min=4,max=100"`
	Password string `json:"password" validate:"required,is-strong-password"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyAccountRequest struct {
	Token string `validate:"required,max=100"`
}
