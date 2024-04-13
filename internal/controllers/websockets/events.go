package websockets

import (
	"github.com/gorilla/websocket"
	"itrevolution-backend/internal/domain"
	"itrevolution-backend/internal/types"
)

func EventWebSocketHandler(serverCtx types.ServerContext, connection *websocket.Conn, user domain.User, stateChan chan int) {
	if _, ok := serverCtx.WsConns[user.ID]; ok {
		exists := false
		for _, wsConn := range serverCtx.WsConns[user.ID] {
			if wsConn == connection {
				exists = true
				break
			}
		}
		if !exists {
			serverCtx.WsConns[user.ID] = append(serverCtx.WsConns[user.ID], connection)
		}
	} else {
		serverCtx.WsConns[user.ID] = make([]*websocket.Conn, 0)
		serverCtx.WsConns[user.ID] = append(serverCtx.WsConns[user.ID], connection)
	}
	stateChan <- WEBSOCKET_STATE_CONTINUE
}
