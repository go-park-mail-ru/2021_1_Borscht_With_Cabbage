package usecase

import (
	"context"
	"sort"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type chatUsecase struct {
	ChatRepo chat.ChatRepo
}

func NewChatUsecase(chatRepo chat.ChatRepo) chat.ChatUsecase {
	return &chatUsecase{
		ChatRepo: chatRepo,
	}
}

func (ch chatUsecase) GetAllChats(ctx context.Context) ([]models.BriefInfoChat, error) {
	if user, ok := checkUser(ctx); ok {
		return ch.ChatRepo.GetAllChatsUser(ctx, user.Uid)
	}

	if restaurant, ok := checkRestaurant(ctx); ok {
		return ch.ChatRepo.GetAllChatsRestaurant(ctx, restaurant.ID)
	}

	foundError := errors.AuthorizationError("not user and restaurant")
	logger.UsecaseLevel().ErrorLog(ctx, foundError)
	return nil, foundError
}

func (ch chatUsecase) GetAllMessages(ctx context.Context, id int) (
	result *models.InfoChat, err error) {

	if user, ok := checkUser(ctx); ok {
		result, err = ch.getAllMsgUser(ctx, *user, id)
	}

	if restaurant, ok := checkRestaurant(ctx); ok {
		result, err = ch.getAllMsgRest(ctx, *restaurant, id)
	}

	if err != nil {
		return nil, err
	}

	sort.SliceStable(result.Messages, func(i, j int) bool {
		return result.Messages[i].Id > result.Messages[j].Id
	})
	return result, nil
}

func (ch chatUsecase) getAllMsgUser(ctx context.Context, user models.User, id int) (
	*models.InfoChat, error) {

	result := new(models.InfoChat)

	infoOpponent, err := ch.ChatRepo.GetRestaurant(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"opponent": infoOpponent})
	result.InfoOpponent = *infoOpponent

	fromMe, err := ch.ChatRepo.GetAllMessagesFromUser(ctx, user.Uid, id)
	if err != nil {
		return nil, err
	}
	for i := range fromMe {
		fromMe[i].FromMe = true
	}

	toMe, err := ch.ChatRepo.GetAllMessagesFromRestaurant(ctx, id, user.Uid)
	if err != nil {
		return nil, err
	}
	for i := range toMe {
		toMe[i].FromMe = false
	}

	result.Messages = append(fromMe, toMe...)
	return result, nil
}

func (ch chatUsecase) getAllMsgRest(ctx context.Context, restaurant models.RestaurantInfo, id int) (
	*models.InfoChat, error) {

	result := new(models.InfoChat)

	infoOpponent, err := ch.ChatRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	result.InfoOpponent = *infoOpponent
	fromMe, err := ch.ChatRepo.GetAllMessagesFromRestaurant(ctx, restaurant.ID, id)
	if err != nil {
		return nil, err
	}
	for i := range fromMe {
		fromMe[i].FromMe = true
	}

	toMe, err := ch.ChatRepo.GetAllMessagesFromRestaurant(ctx, id, restaurant.ID)
	if err != nil {
		return nil, err
	}
	for i := range toMe {
		toMe[i].FromMe = false
	}

	result.Messages = append(fromMe, toMe...)
	return result, nil
}
