package chat

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler interface {
	Connect(c echo.Context) error
	GetKey(c echo.Context) error
}

type WebSocketUsecase interface {
	Connect(ctx context.Context, ws *websocket.Conn) error
	UnConnect(ctx context.Context, ws *websocket.Conn) error
	MessageCame(ctx context.Context, ws *websocket.Conn, msg models.FromClient) error
}

type WebSocketRepo interface {
	SaveMessageFromUser(ctx context.Context, info models.WsMessageForRepo) (mid int, err error)
	SaveMessageFromRestaurant(ctx context.Context, info models.WsMessageForRepo) (mid int, err error)
}

type ChatHandler interface {
	GetAllChats(c echo.Context) error
}

type ChatUsecase interface {
	GetAllChats(ctx context.Context) ([]models.BriefInfoChat, error)
}

type ChatRepo interface {
	GetAllChatsUser(ctx context.Context, uid int) ([]models.BriefInfoChat, error)
	GetAllChatsRestaurant(ctx context.Context, rid int) ([]models.BriefInfoChat, error)
}
