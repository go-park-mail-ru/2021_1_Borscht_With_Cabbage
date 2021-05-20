package notifications

import (
	"errors"
	"github.com/borscht/backend/utils/websocketPool"
)

type OrderNotificator interface {
	OrderStatusChangedNotification(status string, uid int) error
}

type orderNotificator struct {
	UserConnectionsPool       *websocketPool.ConnectionPool
	RestaurantConnectionsPool *websocketPool.ConnectionPool
}

type StatusMessage struct {
	Uid       int    `json:"uid"`
	NewStatus string `json:"status"`
}

func NewOrderNotificator(userPool, restaurantPool *websocketPool.ConnectionPool) OrderNotificator {
	return &orderNotificator{
		UserConnectionsPool:       userPool,
		RestaurantConnectionsPool: restaurantPool,
	}
}

func (N orderNotificator) OrderStatusChangedNotification(status string, uid int) error {
	wsSent := N.UserConnectionsPool.Connections[uid]
	notification := StatusMessage{
		Uid:       uid,
		NewStatus: status,
	}
	if wsSent != nil {
		return wsSent.WriteJSON(notification)
	}

	return errors.New("websocket connection error")
}
