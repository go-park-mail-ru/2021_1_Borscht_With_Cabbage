package usecase

import (
	"context"
	"sync"

	"github.com/borscht/backend/internal/models"
	custWebsocket "github.com/borscht/backend/internal/websocket"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
)

var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

type WebSocketUsecase struct {
}

func NewWebSocketUsecase() custWebsocket.WebSocketUsecase {
	return &WebSocketUsecase{}
}

func (w WebSocketUsecase) Connect(ctx context.Context, ws *websocket.Conn) error {
	connectionPool.Lock()
	connectionPool.connections[ws] = struct{}{}
	connectionPool.Unlock()

	return nil
}

func (w WebSocketUsecase) UnConnect(ctx context.Context, ws *websocket.Conn) error {
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Close websocket": ws.RemoteAddr()})
	connectionPool.Lock()
	delete(connectionPool.connections, ws)
	connectionPool.Unlock()

	return nil
}

func (w WebSocketUsecase) MessageCame(ctx context.Context, ws *websocket.Conn, msg string) error {
	return sendMessageToAllPool(ctx, msg)
}

func connectUser(ctx context.Context, user models.User) error {
	return nil
}

func connectRestaurant(ctx context.Context, user models.RestaurantInfo) error {
	return nil
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
