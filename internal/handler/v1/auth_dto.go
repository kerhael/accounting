package v1

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`         // bearer token
	RefreshToken string `json:"refresh_token"` // refresh token
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`         // bearer token
	RefreshToken string `json:"refresh_token"` // refresh token
}
