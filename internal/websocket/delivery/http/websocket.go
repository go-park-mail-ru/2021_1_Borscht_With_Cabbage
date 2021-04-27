package http

import (
	"context"
	"net/http"
	"sync"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/session"
	custWebsocket "github.com/borscht/backend/internal/websocket"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	WsUsecase      custWebsocket.WebSocketUsecase
	SessionUsecase session.SessionUsecase
}

func NewWebSocketHandler(wsUsecase custWebsocket.WebSocketUsecase,
	sessionUsecase session.SessionUsecase) custWebsocket.WebSocketHandler {
	return &WebSocketHandler{
		WsUsecase:      wsUsecase,
		SessionUsecase: sessionUsecase,
	}
}

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (w WebSocketHandler) GetKey(c echo.Context) error {
	ctx := models.GetContext(c)

	resuslt, err := w.SessionUsecase.CreateKey(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	return models.SendResponse(c, &models.Key{Key: resuslt})
}

func (w WebSocketHandler) Connect(c echo.Context) error {
	ctx := models.GetContext(c)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, failError)
		return failError
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Open websocket": ws.RemoteAddr()})

	connectionPool.Lock()
	connectionPool.connections[ws] = struct{}{}

	defer func(connection *websocket.Conn) {
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Close websocket": ws.LocalAddr()})
		connectionPool.Lock()
		delete(connectionPool.connections, connection)
		connectionPool.Unlock()
	}(ws)

	connectionPool.Unlock()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"error": failError.Error()})
			return failError
		}
		err = sendMessageToAllPool(ctx, string(msg))
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, failError)
			return failError
		}
		logger.DeliveryLevel().DebugLog(ctx, logger.Fields{"message": string(msg)})
	}

	return err
}

func sendMessageToAllPool(ctx context.Context, message string) error {
	connectionPool.RLock()
	defer connectionPool.RUnlock()
	for connection := range connectionPool.connections {
		err := connection.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, failError)
			return err
		}
	}
	return nil
}
