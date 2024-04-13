package auth

import (
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	var req loginRequest
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
	if err := server.DB.First(&user, "name = ?", req.Username).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("invalid username or password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("invalid username or password"))
		return
	}

	tokens, err := generateTokens(user, server.Config)
	if err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, RefreshTokenResponse{
		Response: types.Response{},
		Tokens:   tokens,
	})

}
