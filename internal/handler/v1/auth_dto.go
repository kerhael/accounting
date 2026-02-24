package v1

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"` // bearer token
	User  UserResponse `json:"user"`  // user object
}
