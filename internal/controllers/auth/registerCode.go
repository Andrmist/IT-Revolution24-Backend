package auth

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

type registerCodeRequest struct {
	Code int `json:"code" validate:"required"`
}

func RegisterCode(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	var req registerCodeRequest
	err := render.DecodeJSON(r.Body, &req)

	if err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to decode body"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		domain.HTTPError(w, r, http.StatusBadRequest, domain.ValidationError(validateErr))
		return
	}

	var user domain.User
	if err := server.DB.First(&user, "auth_code = ?", req.Code).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("invalid auth code"))
		return
	}

	user.AuthCode = -1
	user.IsRegistered = true
	if err := server.DB.Save(&user).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	responseTokens, err := generateTokens(user, server.Config)

	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, RefreshTokenResponse{
		Response: types.Response{},
		Tokens:   responseTokens,
		User:     user,
	})
}
