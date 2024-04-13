package users

import (
	"encoding/json"
	"itrevolution-backend/internal/domain"
	"net/http"
)

func UsersGetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(domain.User)
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}
