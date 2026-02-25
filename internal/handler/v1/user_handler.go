package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kerhael/accounting/internal/auth"
	"github.com/kerhael/accounting/internal/domain"
	"github.com/kerhael/accounting/internal/handler/utils"
	"github.com/kerhael/accounting/internal/service"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

// Create a user
// @Summary      Create a user
// @Description Create a new user. A rate limiter prevents from brute force attacks (speed 1s, burst 5)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User payload"
// @Success      201       {object}   UserResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      429       {object}   ErrorResponse  "Too many requests error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/users/ [post]
func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.FirstName == "" {
		http.Error(w, "firstName is required", http.StatusBadRequest)
		return
	}
	if req.LastName == "" {
		http.Error(w, "lastName is required", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toUserResponse(user))
}

// Retrieve authenticated user
// @Summary      Retrieve the authenticated user
// @Description Retrieve the authenticated user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200       {object}   UserResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      404       {object}   ErrorResponse  "User not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Router       /api/v1/users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	user, err := h.service.FindById(r.Context(), userID)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			utils.WriteJSONError(w, http.StatusNotFound, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
