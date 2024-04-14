package messages

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

type readMessageResponse struct {
	types.Response
	Ok bool `json:"ok"`
}

func ReadMessage(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value("server").(types.ServerContext)
	if msgId := chi.URLParam(r, "id"); msgId != "" {
		var msg domain.Message
		if err := server.DB.First(&msg, msgId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			domain.HTTPError(w, r, http.StatusNotFound, errors.New("message id not found"))
			return
		}
		msg.IsRead = true

		if err := server.DB.Save(&msg).Error; err != nil {
			domain.HTTPInternalServerError(w, r, err)
			return
		}
		render.JSON(w, r, readMessageResponse{Ok: true})
	} else {
		domain.HTTPError(w, r, http.StatusBadRequest, errors.New("message id is missing"))
		return
	}
}
