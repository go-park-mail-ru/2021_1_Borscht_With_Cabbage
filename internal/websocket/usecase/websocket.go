package usecase

import (
	custWebsocket "github.com/borscht/backend/internal/websocket"
)

type WebSocketUsecase struct {
}

func NewWebSocketUsecase() custWebsocket.WebSocketUsecase {
	return &WebSocketUsecase{}
}
