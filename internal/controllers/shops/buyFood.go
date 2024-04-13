package shops

import (
	"encoding/json"
	"errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type BuyFoodRequest struct {
	Type string `json:"type" validate:"required"`
}

func ShopsBuyFood(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	if user.Role == "children" {
		var req BuyFoodRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to decode body"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			domain.HTTPError(w, r, http.StatusBadRequest, domain.ValidationError(validateErr))
			return
		}

		var shop domain.FoodShop
		if err := serverCtx.DB.Where("type = ?", req.Type).First(&shop).Error; err != nil {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("pet don't exist"))
			return
		}

		if user.Balance < shop.Cost {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("you don't have enough money"))
			return
		}

		user.Balance = user.Balance - shop.Cost

		if err := serverCtx.DB.Where("id = ?", user.ID).Save(&user).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		food := domain.Food{
			Type:   req.Type,
			Cost:   shop.Cost,
			UserID: user.ID,
		}

		if err := serverCtx.DB.Create(&food).Error; err != nil {
			domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to create pet"))
			return
		}

		render.JSON(w, r, "success")
		return
	}

	render.JSON(w, r, types.Response{Error: "parent don't exist"})
}
