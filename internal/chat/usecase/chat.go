package usecase

import (
	"context"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	serviceChat "github.com/borscht/backend/internal/services/chat"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
)

type chatUsecase struct {
	poolUsers      models.ConnectionPool
	poolRestaurant models.ConnectionPool
	ChatRepo       chat.ChatRepo
	ChatService    serviceChat.ServiceChat
}

func NewChatUsecase(chatRepo chat.ChatRepo, chatService serviceChat.ServiceChat) chat.ChatUsecase {
	return &chatUsecase{
		ChatRepo:       chatRepo,
		ChatService:    chatService,
		poolUsers:      models.NewConnectionPool(),
		poolRestaurant: models.NewConnectionPool(),
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

func (ch *chatUsecase) Connect(ctx context.Context, ws *websocket.Conn) error {
	if user, ok := checkUser(ctx); ok {
		ch.poolUsers.Add(user.Uid, ws)
		return nil
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		ch.poolRestaurant.Add(restaurant.ID, ws)
		return nil
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (ch *chatUsecase) UnConnect(ctx context.Context, ws *websocket.Conn) error {
	if user, ok := checkUser(ctx); ok {
		ch.poolUsers.Remove(user.Uid)
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		ch.poolRestaurant.Remove(restaurant.ID)
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (ch *chatUsecase) MessageCame(ctx context.Context, ws *websocket.Conn, msg models.FromClient) error {
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"message": msg})

	saveInfo := models.WsMessageForRepo{
		Date:     msg.Payload.Message.Date,
		Content:  msg.Payload.Message.Text,
		SentToId: msg.Payload.To.Id,
	}
	if user, ok := checkUser(ctx); ok {
		// сохранение в бд
		saveInfo.SentFromId = user.Uid
		mid, err := ch.ChatRepo.SaveMessageFromUser(ctx, saveInfo)
		if err != nil {
			return err
		}

		// отправить пользователю
		msg.Payload.Message.Id = mid
		return ch.MessageFromUser(ctx, *user, ws, msg)
	}
	if restaurant, ok := checkRestaurant(ctx); ok {
		// сохранение в бд
		saveInfo.SentFromId = restaurant.ID
		mid, err := ch.ChatRepo.SaveMessageFromRestaurant(ctx, saveInfo)
		if err != nil {
			return err
		}

		// отпавить пользователю
		msg.Payload.Message.Id = mid
		return ch.MessageFromRestaurant(ctx, *restaurant, ws, msg)
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (ch *chatUsecase) MessageFromUser(ctx context.Context,
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

	sendConnect := ch.poolRestaurant.Get(toOpponent)
	if sendConnect != nil {
		return sendConnect.WriteJSON(msgSend)
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "not id in connection pool")
	return nil
}

func (ch *chatUsecase) MessageFromRestaurant(ctx context.Context,
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

	sendConnect := ch.poolUsers.Get(toOpponent)
	if sendConnect != nil {
		return sendConnect.WriteJSON(msgSend)
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "not id in connection pool")
	return nil
}

func (ch *chatUsecase) GetAllChats(ctx context.Context) ([]models.BriefInfoChat, error) {
	if user, ok := checkUser(ctx); ok {
		return ch.ChatService.GetAllChats(ctx, user.Uid, 0)

	} else if restaurant, ok := checkRestaurant(ctx); ok {
		return ch.ChatService.GetAllChats(ctx, 0, restaurant.ID)
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return nil, foundError
}

func (ch *chatUsecase) GetAllMessages(ctx context.Context, id int) (*models.InfoChat, error) {

	result := new(models.InfoChat)
	if user, ok := checkUser(ctx); ok {
		messages, err := ch.ChatService.GetAllMessagesUser(ctx, user.Uid, id)
		if err != nil {
			return nil, err
		}

		infoOpponent, err := ch.ChatRepo.GetRestaurant(ctx, id)
		if err != nil {
			return nil, err
		}
		result.Messages = messages
		result.InfoOpponent = *infoOpponent

	} else if restaurant, ok := checkRestaurant(ctx); ok {
		messages, err := ch.ChatService.GetAllMessagesRestaurant(ctx, id, restaurant.ID)
		if err != nil {
			return nil, err
		}

		infoOpponent, err := ch.ChatRepo.GetUser(ctx, id)
		if err != nil {
			return nil, err
		}
		result.Messages = messages
		result.InfoOpponent = *infoOpponent
	}

	return result, nil
}
