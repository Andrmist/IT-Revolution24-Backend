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

type FeedPetRequest struct {
	types.Response
	PetID uint `json:"petId" validate:"required"`
}

func PetsFeed(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	if user.Role == "children" {
		var req FeedPetRequest
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
		if err := server.DB.Where("id = ? and user_id = ?", req.PetID, user.ID).First(&pet).Error; err != nil {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to feed pet"))
			return
		}

		var food domain.Food
		if err := server.DB.Where("user_id = ?", pet.UserID).First(&food).Error; err != nil {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("you don't have food"))
			return
		}

		pet.Satiety = pet.Satiety + 25

		if pet.Satiety > 100 {
			pet.Satiety = 100
		}

		if err := server.DB.Where("id = ? and user_id = ?", req.PetID, user.ID).Save(&pet).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		server.DB.Where("user_id = ?", user.ID).Delete(&food).Commit()

		render.JSON(w, r, "success")
		return
	}

	render.JSON(w, r, types.Response{Error: "parent don't exist"})
}
