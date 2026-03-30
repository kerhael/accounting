package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
// @Router       /users/ [post]
func (h *UserHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.FirstName == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "firstName is required")
		return
	}
	if req.LastName == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "lastName is required")
		return
	}
	if req.Email == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "email is required")
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "password is required")
		return
	}
	if len(req.Password) < 8 {
		utils.WriteJSONError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	user, err := h.service.Create(r.Context(), req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, toUserResponse(user))
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
// @Security BearerAuth
// @Router       /users/me [get]
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
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}

// Update a user
// @Summary      Update a user
// @Description  Update a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 id path int true "User ID"
// @Param        user  body      PatchUserByIdRequest  true  "User payload"
// @Success      200       {object}   UserResponse
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      404       {object}   ErrorResponse  "Not found error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /users/{id} [patch]
func (h *UserHandler) PatchUserById(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if userId != id {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user cannot be updated")
		return
	}

	var req PatchUserByIdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	firstName := ""
	if req.FirstName != nil {
		cleanName := strings.TrimSpace(*req.FirstName)
		firstName = cleanName
	}

	lastName := ""
	if req.LastName != nil {
		cleanName := strings.TrimSpace(*req.LastName)
		lastName = cleanName
	}

	password := ""
	if req.Password != nil {
		if len(*req.Password) < 8 {
			utils.WriteJSONError(w, http.StatusBadRequest, "password must be at least 8 characters")
			return
		}
		password = *req.Password
	}

	user, err := h.service.PatchById(r.Context(), userId, firstName, lastName, password)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		if error, ok := errors.AsType[*domain.EntityNotFoundError](err); ok {
			utils.WriteJSONError(w, http.StatusNotFound, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserResponse(user))
}

// Delete a user
// @Summary      Delete a user
// @Description Delete a user by id
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		id path int true "User ID"
// @Success      204       "No Content"
// @Failure      400       {object}   ErrorResponse  "Bad request error"
// @Failure      401       {object}   ErrorResponse  "Unauthorized error"
// @Failure      500       {object}   ErrorResponse  "Internal server error"
// @Security BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	_, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.service.DeleteById(r.Context(), id)
	if err != nil {
		if error, ok := errors.AsType[*domain.InvalidEntityError](err); ok {
			utils.WriteJSONError(w, http.StatusBadRequest, error.Error())
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
