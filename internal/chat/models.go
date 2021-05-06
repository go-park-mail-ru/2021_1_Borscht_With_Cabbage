package chat

import (
	"context"

	"github.com/borscht/backend/internal/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ChatHandler interface {
	Connect(c echo.Context) error
	GetKey(c echo.Context) error
	GetAllChats(c echo.Context) error
	GetAllMessages(c echo.Context) error
}

type ChatUsecase interface {
	GetAllChats(ctx context.Context) ([]models.BriefInfoChat, error)
	GetAllMessages(ctx context.Context, id int) (*models.InfoChat, error)
	Connect(ctx context.Context, ws *websocket.Conn) error
	UnConnect(ctx context.Context, ws *websocket.Conn) error
	MessageCame(ctx context.Context, ws *websocket.Conn, msg models.FromClient) error
}
