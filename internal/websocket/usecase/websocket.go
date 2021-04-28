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

type connectionPool struct {
	sync.RWMutex
	connections map[int]*websocket.Conn
}

type WebSocketUsecase struct {
	poolUsers      connectionPool
	poolRestaurant connectionPool
	WsRepo         custWebsocket.WebSocketRepo
}

func NewWebSocketUsecase(wsRepo custWebsocket.WebSocketRepo) custWebsocket.WebSocketUsecase {
	return &WebSocketUsecase{
		WsRepo: wsRepo,
		poolUsers: connectionPool{
			connections: make(map[int]*websocket.Conn),
		},
		poolRestaurant: connectionPool{
			connections: make(map[int]*websocket.Conn),
		},
	}
}

func checkUser(ctx context.Context) (*models.User, bool) {
	userInterfece := ctx.Value("User")
	if userInterfece == nil {
		logger.UsecaseLevel().InlineDebugLog(ctx, "not user in context")
		return nil, false
	}

	user, ok := userInterfece.(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, false
	}

	return &user, true
}

func checkRestaurant(ctx context.Context) (*models.RestaurantInfo, bool) {
	restaurantInterfece := ctx.Value("Restaurant")
	if restaurantInterfece == nil {
		logger.UsecaseLevel().InlineDebugLog(ctx, "not restaurant in context")
		return nil, false
	}

	restaurant, ok := restaurantInterfece.(models.RestaurantInfo)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.Restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, false
	}

	return &restaurant, true
}

func (w *WebSocketUsecase) Connect(ctx context.Context, ws *websocket.Conn) error {
	if user, ok := checkUser(ctx); ok {
		w.poolUsers.Lock()
		w.poolUsers.connections[user.Uid] = ws
		w.poolUsers.Unlock()

		return nil
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		w.poolRestaurant.Lock()
		w.poolRestaurant.connections[restaurant.ID] = ws
		w.poolRestaurant.Unlock()

		return nil
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (w *WebSocketUsecase) UnConnect(ctx context.Context, ws *websocket.Conn) error {
	if user, ok := checkUser(ctx); ok {
		w.poolUsers.Lock()
		delete(w.poolUsers.connections, user.Uid)
		w.poolUsers.Unlock()
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		w.poolRestaurant.Lock()
		delete(w.poolRestaurant.connections, restaurant.ID)
		w.poolRestaurant.Unlock()
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (w *WebSocketUsecase) MessageCame(ctx context.Context, ws *websocket.Conn, msg models.FromClient) error {
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"message": msg})

	saveInfo := models.WsMessageForRepo{
		Date:     msg.Payload.Message.Date,
		Content:  msg.Payload.Message.Text,
		SentToId: msg.Payload.To.Id,
	}
	if user, ok := checkUser(ctx); ok {
		// сохранение в бд
		saveInfo.SentFromId = user.Uid
		mid, err := w.WsRepo.SaveMessageFromUser(ctx, saveInfo)
		if err != nil {
			return err
		}

		// отправить пользователю
		msg.Payload.Message.Id = mid
		return w.MessageFromUser(ctx, *user, ws, msg)
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		// сохранение в бд
		saveInfo.SentFromId = restaurant.ID
		mid, err := w.WsRepo.SaveMessageFromRestaurant(ctx, saveInfo)
		if err != nil {
			return err
		}

		// отпавить пользователю
		msg.Payload.Message.Id = mid
		return w.MessageFromRestaurant(ctx, *restaurant, ws, msg)
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (w *WebSocketUsecase) MessageFromUser(ctx context.Context,
	user models.User, ws *websocket.Conn, msg models.FromClient) error {

	toOpponent := msg.Payload.To.Id
	msgSend := models.ToClient{
		Action: "message",
	}
	msgSend.Payload.Message = msg.Payload.Message
	msgSend.Payload.From = models.WsOpponent{
		Id:     user.Uid,
		Avatar: user.Avatar,
		Name:   user.Name,
	}

	sendConnect := w.poolRestaurant.connections[toOpponent]
	if sendConnect != nil {
		return sendConnect.WriteJSON(msgSend)
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "not id in connection pool")
	return nil
}

func (w *WebSocketUsecase) MessageFromRestaurant(ctx context.Context,
	restaurant models.RestaurantInfo, ws *websocket.Conn, msg models.FromClient) error {

	toOpponent := msg.Payload.To.Id
	msgSend := models.ToClient{
		Action: "message",
	}
	msgSend.Payload.Message = msg.Payload.Message
	msgSend.Payload.From = models.WsOpponent{
		Id:     restaurant.ID,
		Avatar: restaurant.Avatar,
		Name:   restaurant.Title,
	}

	sendConnect := w.poolUsers.connections[toOpponent]
	if sendConnect != nil {
		return sendConnect.WriteJSON(msgSend)
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "not id in connection pool")
	return nil
}
