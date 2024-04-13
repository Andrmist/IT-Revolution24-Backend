package pets

import (
	"errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
)

func PetsGetAll(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var pets []domain.Pet

	if err := server.DB.Where("user_id = ?", user.ID).Find(&pets).Error; err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to get all pets"))
	}

	render.JSON(w, r, pets)
}
