package pets

import (
	"errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
)

type getAllPetsResponse struct {
	ID      uint   `json:"id"`
	Type    string `json:"type"`
	Sex     string `json:"sex"`
	Satiety int    `json:"satiety"`
}

func PetsGetAll(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var pets []domain.Pet

	if err := server.DB.Where("user_id = ?", user.ID).Find(&pets).Error; err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to get all pets"))
	}

	var response []getAllPetsResponse

	for _, pet := range pets {
		response = append(response, getAllPetsResponse{
			ID:      pet.ID,
			Type:    pet.Type,
			Sex:     pet.Sex,
			Satiety: pet.Satiety,
		})
	}

	render.JSON(w, r, response)
}
