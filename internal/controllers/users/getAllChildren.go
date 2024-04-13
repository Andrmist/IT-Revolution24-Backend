package users

import (
	"encoding/json"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

func UsersGetChildren(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	user := r.Context().Value("user").(domain.User)

	var childrens []domain.User
	if err := serverCtx.DB.Where("email = ?", user.Email).Find(&childrens).Error; err != nil {
		domain.HTTPInternalServerError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(childrens)
}
