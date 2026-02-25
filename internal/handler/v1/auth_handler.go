package v1

import (
	"encoding/json"
	"net/http"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/service"
	"github.com/kerhael/accounting/pkg/security"
)

type AuthHandler struct {
	userService service.UserServiceInterface
	jwtService  *auth.JWTService
}

func NewAuthHandler(userService service.UserServiceInterface, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Login
// @Summary      Login
// @Description User login. A rate limiter prevents from brute force attacks (speed 1s, burst 5)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body      LoginRequest  true  "Login payload"
// @Success      200       {object}   LoginResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      429       {object}   ErrorResponse  "Too many requests error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/users/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	err = security.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	token, err := h.jwtService.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "could not generate token", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
	})
}
