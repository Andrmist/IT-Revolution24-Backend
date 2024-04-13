package users

import (
	"encoding/json"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

type getAllChildrenResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float32 `json:"balance"`

	AlivePetsCount    int64 `json:"alivePetsCount"`
	StarvingPetsCount int64 `json:"starvingPetsCount"`

	Messages []string `json:"messages"`
}

func UsersGetChildren(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var childrens []domain.User
	if err := serverCtx.DB.Where("email = ? and role = ?", user.Email, "children").Find(&childrens).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	var response []getAllChildrenResponse
	for _, children := range childrens {
		var alivePetsCount, starvingPetsCount int64

		if err := serverCtx.DB.Model(&domain.Pet{}).Where("user_id = ?", children.ID).Count(&alivePetsCount).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		if err := serverCtx.DB.Model(&domain.Pet{}).Where("user_id = ?", children.ID).Count(&starvingPetsCount).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}

		response = append(response, getAllChildrenResponse{
			ID:       children.ID,
			Name:     children.Name,
			Email:    children.Email,
			Password: children.Password,
			Balance:  children.Balance,

			AlivePetsCount:    alivePetsCount,
			StarvingPetsCount: starvingPetsCount,

			Messages: make([]string, 0),
		})
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
