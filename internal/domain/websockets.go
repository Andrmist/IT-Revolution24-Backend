package domain

type WebSocketMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
