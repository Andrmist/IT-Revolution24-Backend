package websockets

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"itrevolution-backend/internal/controllers/auth"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	WEBSOCKET_STATE_CONTINUE = 1
	WEBSOCKET_STATE_CLOSE    = 0
)

func WebSocketAuthHandler(eventHandler func(serverCtx types.ServerContext, connection *websocket.Conn, user domain.User, stateChan chan int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverCtx := r.Context().Value("server").(types.ServerContext)
		connection, _ := upgrader.Upgrade(w, r, nil)

		var user domain.User
		stateChan := make(chan int)

		go func() {
			for {
				mt, message, err := connection.ReadMessage()

				if err != nil || mt == websocket.CloseMessage {
					break
				}

				var parsedMessage types.WebSocketMessage
				if err := json.Unmarshal(message, &parsedMessage); err != nil {
					serverCtx.Log.Error(err)
					continue
				}
				if parsedMessage.Event == "auth" {
					if accessToken, ok := parsedMessage.Data.(string); ok {
						user, err = auth.GetUserFromAccessToken(serverCtx, accessToken)
						if err != nil {
							msg, _ := json.Marshal(types.WebSocketMessage{
								Event: "error",
								Data:  err,
							})
							connection.WriteMessage(websocket.TextMessage, msg)
							continue
						}
						msg, _ := json.Marshal(types.WebSocketMessage{
							Event: "auth",
							Data:  "ok",
						})
						connection.WriteMessage(websocket.TextMessage, msg)
						continue
					}
				}
			}
		}()
		for {
			if user.ID != 0 {
				go eventHandler(serverCtx, connection, user, stateChan)
				switch <-stateChan {
				case WEBSOCKET_STATE_CONTINUE:
					continue
				case WEBSOCKET_STATE_CLOSE:
					break
				}
			}
		}

		defer connection.Close()
	}
}
