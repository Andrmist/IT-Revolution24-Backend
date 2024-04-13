package auth

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	types.Response
}

func Login(w http.ResponseWriter, r *http.Request) {
	//server := r.Context().Value("server").(types.ServerContext)
	var req loginRequest
	err := render.DecodeJSON(r.Body, &req)

	if err != nil {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("failed to decode body"))
		return
	}

}
