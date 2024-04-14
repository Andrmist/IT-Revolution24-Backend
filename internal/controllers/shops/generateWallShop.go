package shops

import (
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
)

type shopsWallResponse struct {
	types.Response

	Foods []domain.FoodShop `json:"foods"`
	Pets  []domain.PetShop  `json:"pets"`
}

func ShopsGenerateWallShop(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var response shopsWallResponse

	if user.Role == "child" {
		if err := serverCtx.DB.Find(&response.Foods).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := serverCtx.DB.Find(&response.Pets).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		render.JSON(w, r, response)
		return
	}

	render.JSON(w, r, types.Response{Error: "parent don't exist"})
}
