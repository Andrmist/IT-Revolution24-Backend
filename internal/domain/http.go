package domain

import (
	"fmt"
	"github.com/go-chi/render"
	"itrevolution-backend/internal/types"
	"net/http"
)

func HTTPInternalServerError(w http.ResponseWriter, r *http.Request, error error) {
	server := r.Context().Value("server").(types.ServerContext)
	server.Log.Error(error)
	w.WriteHeader(http.StatusInternalServerError)
	render.JSON(w, r, types.Response{Error: error.Error()})
}

func HTTPError(w http.ResponseWriter, r *http.Request, status int, error error) {
	w.WriteHeader(status)
	if error != nil {
		render.JSON(w, r, types.Response{Error: error.Error()})
	} else {
		errorResponse := ""
		switch status {
		case http.StatusUnauthorized:
			errorResponse = "Unathorized"
		case http.StatusForbidden:
			errorResponse = "Forbidden"
		case http.StatusBadRequest:
			errorResponse = "BadRequest"
		}
		if errorResponse != "" {
			render.JSON(w, r, types.Response{Error: errorResponse})
		} else {
			render.JSON(w, r, types.Response{Error: fmt.Sprintf("%v", status)})
		}
	}
}
