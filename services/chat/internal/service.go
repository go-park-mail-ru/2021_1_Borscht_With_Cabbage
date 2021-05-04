package internal

import (
	"context"
	"sort"

	chatRepo "github.com/borscht/backend/services/chat/repository"
	protoChat "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type service struct {
	chatRepo chatRepo.ChatRepo
}

func NewService(chatRepo chatRepo.ChatRepo) *service {
	return &service{
		chatRepo: chatRepo,
	}
}

// func checkUser(ctx context.Context) (*models.User, bool) {
// 	userInterfece := ctx.Value("User")
// 	if userInterfece == nil {
// 		logger.UsecaseLevel().InlineDebugLog(ctx, "not user in context")
// 		return nil, false
// 	}

// 	user, ok := userInterfece.(models.User)
// 	if !ok {
// 		failError := errors.FailServerError("failed to convert to models.User")
// 		logger.UsecaseLevel().ErrorLog(ctx, failError)
// 		return nil, false
// 	}

// 	return &user, true
// }

// func checkRestaurant(ctx context.Context) (*models.RestaurantInfo, bool) {
// 	restaurantInterfece := ctx.Value("Restaurant")
// 	if restaurantInterfece == nil {
// 		logger.UsecaseLevel().InlineDebugLog(ctx, "not restaurant in context")
// 		return nil, false
// 	}

// 	restaurant, ok := restaurantInterfece.(models.RestaurantInfo)
// 	if !ok {
// 		failError := errors.FailServerError("failed to convert to models.Restaurant")
// 		logger.UsecaseLevel().ErrorLog(ctx, failError)
// 		return nil, false
// 	}

// 	return &restaurant, true
// }

func (s service) GetAllChats(ctx context.Context, Id *protoChat.Id) (*protoChat.MoreBriefInfoChat, error) {
	logger.UsecaseLevel().InfoLog(ctx, logger.Fields{"Uid": Id.Uid, "Rid": Id.Rid})
	if Id.Uid != 0 {
		result, err := s.chatRepo.GetAllChatsUser(ctx, int(Id.Uid))
		if err != nil {
			return nil, err
		}
		return &protoChat.MoreBriefInfoChat{
			More: result,
		}, nil
	} else if Id.Rid != 0 {
		result, err := s.chatRepo.GetAllChatsRestaurant(ctx, int(Id.Rid))
		if err != nil {
			return nil, err
		}
		return &protoChat.MoreBriefInfoChat{
			More: result,
		}, nil
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return nil, foundError
}

func (ch service) GetAllMessagesUser(ctx context.Context, Id *protoChat.Id) (
	result *protoChat.MoreInfoMessage, err error) {

	result = new(protoChat.MoreInfoMessage)

	fromMe, err := ch.chatRepo.GetAllMessagesFromUser(ctx, int(Id.Uid), int(Id.Rid))
	if err != nil {
		return nil, err
	}
	for i := range fromMe {
		fromMe[i].FromMe = true
	}

	toMe, err := ch.chatRepo.GetAllMessagesFromRestaurant(ctx, int(Id.Rid), int(Id.Uid))
	if err != nil {
		return nil, err
	}
	for i := range toMe {
		toMe[i].FromMe = false
	}

	result.More = append(fromMe, toMe...)

	sort.SliceStable(result.More, func(i, j int) bool {
		return result.More[i].Id > result.More[j].Id
	})
	return result, nil
}

func (ch service) GetAllMessagesRestaurant(ctx context.Context, Id *protoChat.Id) (
	result *protoChat.MoreInfoMessage, err error) {

	result = new(protoChat.MoreInfoMessage)

	fromMe, err := ch.chatRepo.GetAllMessagesFromRestaurant(ctx, int(Id.Rid), int(Id.Uid))
	if err != nil {
		return nil, err
	}
	for i := range fromMe {
		fromMe[i].FromMe = true
	}

	toMe, err := ch.chatRepo.GetAllMessagesFromUser(ctx, int(Id.Uid), int(Id.Rid))
	if err != nil {
		return nil, err
	}
	for i := range toMe {
		toMe[i].FromMe = false
	}

	result.More = append(fromMe, toMe...)

	sort.SliceStable(result.More, func(i, j int) bool {
		return result.More[i].Id > result.More[j].Id
	})
	return result, nil
}

// func (ch service) GetAllMessages(ctx context.Context, id *protoChat.Id) (
// 	result *protoChat.InfoChat, err error) {

// 	if user, ok := checkUser(ctx); ok {
// 		result, err = ch.getAllMsgUser(ctx, *user, int(id.Id))
// 	}

// 	if restaurant, ok := checkRestaurant(ctx); ok {
// 		result, err = ch.getAllMsgRest(ctx, *restaurant, int(id.Id))
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	sort.SliceStable(result.Messages, func(i, j int) bool {
// 		return result.Messages[i].Id > result.Messages[j].Id
// 	})
// 	return result, nil
// }

// func (ch service) getAllMsgUser(ctx context.Context, user models.User, id int) (
// 	*protoChat.InfoChat, error) {

// 	result := new(protoChat.InfoChat)

// 	infoOpponent, err := ch.chatRepo.GetRestaurant(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"opponent": infoOpponent})
// 	result.Info = infoOpponent

// 	fromMe, err := ch.chatRepo.GetAllMessagesFromUser(ctx, user.Uid, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range fromMe {
// 		fromMe[i].FromMe = true
// 	}

// 	toMe, err := ch.chatRepo.GetAllMessagesFromRestaurant(ctx, id, user.Uid)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range toMe {
// 		toMe[i].FromMe = false
// 	}

// 	result.Messages = append(fromMe, toMe...)
// 	return result, nil
// }

// func (ch service) getAllMsgRest(ctx context.Context, restaurant models.RestaurantInfo, id int) (
// 	*protoChat.InfoChat, error) {

// 	result := new(protoChat.InfoChat)

// 	infoOpponent, err := ch.chatRepo.GetUser(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result.Info = infoOpponent
// 	fromMe, err := ch.chatRepo.GetAllMessagesFromRestaurant(ctx, restaurant.ID, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range fromMe {
// 		fromMe[i].FromMe = true
// 	}

// 	toMe, err := ch.chatRepo.GetAllMessagesFromUser(ctx, id, restaurant.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for i := range toMe {
// 		toMe[i].FromMe = false
// 	}

// 	result.Messages = append(fromMe, toMe...)
// 	return result, nil
// }

// func getContext() context.Context {
// 	ctx := context.Background()
// 	requestID := fmt.Sprintf("%016x", rand.Int())[:10]
// 	ctx = context.WithValue(ctx, "request_id", requestID)
// 	return ctx
// }

// func (s service) RouteChat(inStream protoChat.Chat_RouteChatServer) error {
// 	for {
// 		msg, err := inStream.Recv()
// 		ctx := getContext()
// 		if err == io.EOF {
// 			logger.UsecaseLevel().ErrorLog(ctx, err)
// 			return nil
// 		}
// 		if err != nil {
// 			logger.UsecaseLevel().ErrorLog(ctx, err)
// 			return err
// 		}

// 		saveInfo := modelsService.WsMessageForRepo{
// 			Date:       msg.Message.Date,
// 			Content:    msg.Message.Text,
// 			SentToId:   msg.To.Id,
// 			SentFromId: msg.From.Id,
// 		}

// 		// сохранение в бд
// 		var mid int
// 		if msg.From.IsUser {
// 			mid, err = s.chatRepo.SaveMessageFromUser(ctx, saveInfo)
// 			if err != nil {
// 				return err
// 			}

// 		} else {
// 			mid, err = s.chatRepo.SaveMessageFromRestaurant(ctx, saveInfo)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// отправить пользователю
// 		msg.Message.Id = int32(mid)
// 		response := protoChat.InfoMessage{
// 			Message: msg.Message,
// 			From:    msg,
// 		}
// 		logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"->": response})
// 		inStream.Send(response)
// 	}
// 	return nil
// }

// func createResponseFromUser(user models.User, msg *protoChat.FromClient) *protoChat.ToClient {

// 	toOpponent := msg.Payload.To.Id
// 	msgSend := protoChat.ToClient{
// 		Action:   "message",
// 		IdTo:     toOpponent,
// 		IsUserTo: false,
// 	}
// 	msgSend.Payload.Message = msg.Payload.Message
// 	msgSend.Payload.From = &protoChat.WsOpponent{
// 		Id:     int32(user.Uid),
// 		Avatar: user.Avatar,
// 		Name:   user.Name,
// 	}

// 	return &msgSend
// }

// func createResponseFromRestaurant(restaurant models.RestaurantInfo,
// 	msg *protoChat.FromClient) *protoChat.ToClient {

// 	toOpponent := msg.Payload.To.Id
// 	msgSend := protoChat.ToClient{
// 		Action:   "message",
// 		IdTo:     toOpponent,
// 		IsUserTo: true,
// 	}
// 	msgSend.Payload.Message = msg.Payload.Message
// 	msgSend.Payload.From = &protoChat.WsOpponent{
// 		Id:     int32(restaurant.ID),
// 		Avatar: restaurant.Avatar,
// 		Name:   restaurant.Title,
// 	}

// 	return &msgSend
// }
