package http

import (
	"fmt"
	"net/http"

	"github.com/borscht/backend/internal/models"
	custWebsocket "github.com/borscht/backend/internal/websocket"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
}

func NewWebSocketHandler() custWebsocket.WebSocketHandler {
	return &WebSocketHandler{}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (w WebSocketHandler) Connect(c echo.Context) error {
	ctx := models.GetContext(c)
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Open websocket": ""})

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Close websocket": ""})
		ws.Close()
	}()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
