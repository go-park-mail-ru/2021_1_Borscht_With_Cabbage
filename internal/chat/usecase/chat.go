package usecase

import (
	"context"

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
