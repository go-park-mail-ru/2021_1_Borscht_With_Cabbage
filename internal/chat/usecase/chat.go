package usecase

import (
	"context"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	serviceAuth "github.com/borscht/backend/internal/services/auth"
	serviceChat "github.com/borscht/backend/internal/services/chat"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
)

type chatUsecase struct {
	poolUsers      models.ConnectionPool
	poolRestaurant models.ConnectionPool
	ChatService    serviceChat.ServiceChat
	AuthService    serviceAuth.ServiceAuth
}

func NewChatUsecase(chatService serviceChat.ServiceChat,
	authService serviceAuth.ServiceAuth) chat.ChatUsecase {
	return &chatUsecase{
		ChatService:    chatService,
		poolUsers:      models.NewConnectionPool(),
		poolRestaurant: models.NewConnectionPool(),
		AuthService:    authService,
	}
}

func checkUser(ctx context.Context) (*models.User, bool) {
	userInContext := ctx.Value("User")
	if userInContext == nil {
		logger.UsecaseLevel().InlineDebugLog(ctx, "not user in context")
		return nil, false
	}

	user, ok := userInContext.(models.User)
	if !ok {
		failError := errors.FailServerError("failed to convert to models.User")
		logger.UsecaseLevel().ErrorLog(ctx, failError)
		return nil, false
	}

	return &user, true
}

func checkRestaurant(ctx context.Context) (*models.RestaurantInfo, bool) {
	restaurantInContext := ctx.Value("Restaurant")
	if restaurantInContext == nil {
		logger.UsecaseLevel().InlineDebugLog(ctx, "not restaurant in context")
		return nil, false
	}

	restaurant, ok := restaurantInContext.(models.RestaurantInfo)
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

	foundError := errors.BadRequestError("not user and restaurant")
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

	foundError := errors.BadRequestError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return foundError
}

func (ch *chatUsecase) ProcessMessage(ctx context.Context, ws *websocket.Conn, msg models.FromClient) error {
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"message": msg})

	message := models.ChatMessage{
		Text: msg.Payload.Message.Text,
		Date: msg.Payload.Message.Date,
	}
	me := new(models.ChatUser)
	opponent := new(models.ChatUser)
	opponent.Id = msg.Payload.To.Id
	wsOpponent := new(models.WsOpponent)
	var wsSent *websocket.Conn

	if user, ok := checkUser(ctx); ok {
		me.Id = user.Uid
		me.Role = config.RoleUser
		opponent.Role = config.RoleAdmin

		wsOpponent.Avatar = user.Avatar
		wsOpponent.Name = user.Name
		wsOpponent.Id = user.Uid
		wsSent = ch.poolRestaurant.Get(opponent.Id)
	} else if restaurant, ok := checkRestaurant(ctx); ok {
		me.Id = restaurant.ID
		me.Role = config.RoleAdmin
		opponent.Role = config.RoleUser

		wsOpponent.Avatar = restaurant.Avatar
		wsOpponent.Name = restaurant.Title
		wsOpponent.Id = restaurant.ID
		wsSent = ch.poolUsers.Get(opponent.Id)
	} else {
		foundError := errors.AuthorizationError("not user and restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, foundError)
		return foundError
	}

	messageSent, err := ch.ChatService.ProcessMessage(ctx, models.InfoChatMessage{
		Message:   message,
		Sender:    *me,
		Recipient: *opponent,
	})

	msgSend := models.ToClient{
		Action: "message",
	}
	msgSend.Payload.Message = msg.Payload.Message
	msgSend.Payload.Message.Id = messageSent.Message.Mid
	msgSend.Payload.From = *wsOpponent

	if wsSent != nil {
		return wsSent.WriteJSON(msgSend)
	}
	logger.UsecaseLevel().InlineInfoLog(ctx, "no id in connection pool")

	return err
}

func (ch *chatUsecase) GetAllChats(ctx context.Context) ([]models.BriefInfoChat, error) {
	var me models.ChatUser
	if user, ok := checkUser(ctx); ok {
		me.Id = user.Uid
		me.Role = config.RoleUser

	} else if restaurant, ok := checkRestaurant(ctx); ok {
		me.Id = restaurant.ID
		me.Role = config.RoleAdmin

	} else {
		foundError := errors.AuthorizationError("not user and restaurant")
		logger.UsecaseLevel().ErrorLog(ctx, foundError)
		return nil, foundError
	}

	messages, err := ch.ChatService.GetAllChats(ctx, me)
	if err != nil {
		return nil, err
	}

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"all chats": messages})
	return ch.convertChats(ctx, me, messages)
}

func (ch *chatUsecase) convertChats(ctx context.Context, me models.ChatUser,
	messages []models.InfoChatMessage) ([]models.BriefInfoChat, error) {

	breifChats := make([]models.BriefInfoChat, 0)
	for _, value := range messages {
		breifChat := new(models.BriefInfoChat)
		breifChat.LastMessage = value.Message.Text

		opponent := value.Sender
		if me == value.Sender {
			opponent = value.Recipient
		}
		breifChat.Uid = opponent.Id

		if opponent.Role == config.RoleAdmin {
			restaurant, err := ch.AuthService.GetByRid(ctx, opponent.Id)
			if err != nil {
				return nil, err
			}
			breifChat.Avatar = restaurant.Avatar
			breifChat.Name = restaurant.Title
		} else {
			user, err := ch.AuthService.GetByUid(ctx, opponent.Id)
			if err != nil {
				return nil, err
			}
			breifChat.Avatar = user.Avatar
			breifChat.Name = user.Name
		}

		breifChats = append(breifChats, *breifChat)
	}

	return breifChats, nil
}

func (ch *chatUsecase) GetAllMessages(ctx context.Context, id int) (*models.InfoChat, error) {
	me := new(models.ChatUser)
	opponent := new(models.ChatUser)
	chat := new(models.InfoChat)
	if user, ok := checkUser(ctx); ok {
		me.Id = user.Uid
		me.Role = config.RoleUser
		opponent.Id = id
		opponent.Role = config.RoleAdmin

		restaurant, err := ch.AuthService.GetByRid(ctx, id)
		if err != nil {
			return nil, err
		}
		chat.Avatar = restaurant.Avatar
		chat.Name = restaurant.Title

	} else if restaurant, ok := checkRestaurant(ctx); ok {
		me.Id = restaurant.ID
		me.Role = config.RoleAdmin
		opponent.Id = id
		opponent.Role = config.RoleUser

		user, err := ch.AuthService.GetByUid(ctx, id)
		if err != nil {
			return nil, err
		}
		chat.Avatar = user.Avatar
		chat.Name = user.Name
	}
	chat.Uid = opponent.Id

	messages, err := ch.ChatService.GetAllMessages(ctx, *me, *opponent)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.convertMessages(ctx, *me, messages)
	if err != nil {
		return nil, err
	}
	chat.Messages = msgs

	return chat, nil
}

func (ch *chatUsecase) convertMessages(ctx context.Context, me models.ChatUser,
	messages []models.InfoChatMessage) ([]models.InfoMessage, error) {

	msgs := make([]models.InfoMessage, 0)
	for _, value := range messages {
		msg := models.InfoMessage{
			Id:   value.Message.Mid,
			Date: value.Message.Date,
			Text: value.Message.Text,
		}

		msg.FromMe = (me == value.Sender)

		msgs = append(msgs, msg)
	}

	return msgs, nil
}
