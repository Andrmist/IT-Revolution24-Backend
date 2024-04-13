package users

import (
	"encoding/json"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUpdateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func UsersUpdateMe(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	userOriginal := r.Context().Value("user").(domain.User)
	user := r.Context().Value("user").(domain.User)

	var userUpdate userUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userOriginal.Role != "admin" {
		user.Role = userOriginal.Role
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), 14)
	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	if err := validator.New().Struct(userUpdate); err != nil {
		validateErr := err.(validator.ValidationErrors)
		domain.HTTPError(w, r, http.StatusBadRequest, domain.ValidationError(validateErr))
		return
	}

	newUser := &domain.User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		DeletedAt:    user.DeletedAt,
		Name:         userUpdate.Name,
		Email:        userUpdate.Email,
		Password:     string(password),
		Role:         user.Role,
		AuthCode:     user.AuthCode,
		IsRegistered: user.IsRegistered,
		Balance:      user.Balance,
		Pets:         user.Pets,
	}

	if err := serverCtx.DB.Save(newUser).Error; err != nil {
		serverCtx.Log.Error(errors.Wrap(err, "failed to updated user"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
