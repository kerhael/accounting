package v1

type CreateUserRequest struct {
	FirstName string `json:"firstName"` // User first name
	LastName  string `json:"lastName"`  // User last name
	Email     string `json:"email"`     // User email
	Password  string `json:"password"`  // User password (minimum 8 characters)
}

type UserResponse struct {
	ID        int    `json:"id"`        // User identifier
	FirstName string `json:"firstName"` // User first name
	LastName  string `json:"lastName"`  // User last name
	Email     string `json:"email"`     // User email
}
