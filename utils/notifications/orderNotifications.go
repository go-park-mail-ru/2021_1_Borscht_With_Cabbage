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

type WebsocketNotification struct {
	Action       string `json:"action"`
	Notification string `json:"notification"`
	NewStatus    string `json:"status"`
}

func NewOrderNotificator(userPool, restaurantPool *websocketPool.ConnectionPool) OrderNotificator {
	return &orderNotificator{
		UserConnectionsPool:       userPool,
		RestaurantConnectionsPool: restaurantPool,
	}
}

func (N orderNotificator) OrderStatusChangedNotification(status string, uid int) error {
	wsSent := N.UserConnectionsPool.Connections[uid]
	notification := WebsocketNotification{
		Action:       "message",
		Notification: "status",
		NewStatus:    status,
	}
	if wsSent != nil {
		return wsSent.WriteJSON(notification)
	}

	return errors.New("websocket connection error")
}
