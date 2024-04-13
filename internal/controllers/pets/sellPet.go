package pets

import (
	"encoding/json"
	"errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type sellPetRequest struct {
	PetId uint `json:"petId" validate:"required"`
}

func PetsSell(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var req sellPetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to decode body"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		domain.HTTPError(w, r, http.StatusBadRequest, domain.ValidationError(validateErr))
		return
	}

	var pet domain.Pet
	if err := server.DB.Where("id = ? and user_id = ?", req.PetId, user.ID).First(&pet).Error; err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to sell pet"))
		return
	}

	user.Balance = user.Balance + pet.Cost

	server.DB.Where("id = ? and user_id = ?", req.PetId, user.ID).Delete(&pet).Commit()

	if err := server.DB.Where("id = ?", user.ID).Save(&user).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	render.JSON(w, r, "sccessfully sold pet")
}
