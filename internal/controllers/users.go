package controllers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

func UsersGetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(domain.User)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UsersGetAll(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	var users domain.User
	if err := serverCtx.DB.Find(&users).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UsersUpdateMe(w http.ResponseWriter, r *http.Request) {
	serverCtx := r.Context().Value("server").(types.ServerContext)
	userOriginal := r.Context().Value("user").(domain.User)
	user := r.Context().Value("user").(domain.User)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userOriginal.Role != "admin" {
		user.Role = userOriginal.Role
	}
	if err := serverCtx.DB.Save(&user).Error; err != nil {
		serverCtx.Log.Error(errors.Wrap(err, "failed to save updated user"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
