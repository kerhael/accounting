package v1

import (
	"encoding/json"
	"net/http"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/handler/utils"
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
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate required fields
	if req.Email == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "email is required")
		return
	}

	if req.Password == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "password is required")
		return
	}

	user, err := h.userService.FindByEmail(r.Context(), req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	err = security.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.jwtService.GenerateJWT(user.ID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "could not generate token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, LoginResponse{
		Token: token,
	})
}
