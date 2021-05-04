package chat

import (
	"context"

	"github.com/borscht/backend/internal/models"
	protoChat "github.com/borscht/backend/services/proto/chat"
)

type ServiceChat interface {
	GetAllChats(ctx context.Context, Uid, Rid int) ([]models.BriefInfoChat, error)
	GetAllMessagesUser(ctx context.Context, Uid, Rid int) ([]models.InfoMessage, error)
	GetAllMessagesRestaurant(ctx context.Context, Uid, Rid int) ([]models.InfoMessage, error)
}

type service struct {
	chatService protoChat.ChatClient
}

func NewService(chatService protoChat.ChatClient) ServiceChat {
	return &service{
		chatService: chatService,
	}
}

func (s service) GetAllChats(ctx context.Context, Uid, Rid int) ([]models.BriefInfoChat, error) {
	messages, err := s.chatService.GetAllChats(ctx, &protoChat.Id{Uid: int32(Uid), Rid: int32(Rid)})
	if err != nil {
		return nil, err
	}

	response := make([]models.BriefInfoChat, 0)
	for _, message := range messages.More {
		opponent := models.InfoOpponent{
			Uid:    int(message.Info.Uid),
			Name:   message.Info.Name,
			Avatar: message.Info.Avatar,
		}
		msg := models.BriefInfoChat{
			InfoOpponent: opponent,
			LastMessage:  message.LastMessage,
		}

		response = append(response, msg)
	}

	return response, nil
}

func (s service) GetAllMessagesUser(ctx context.Context, Uid, Rid int) ([]models.InfoMessage, error) {
	messages, err := s.chatService.GetAllMessagesUser(ctx, &protoChat.Id{Uid: int32(Uid), Rid: int32(Rid)})
	if err != nil {
		return nil, err
	}

	response := make([]models.InfoMessage, 0)
	for _, message := range messages.More {
		msg := models.InfoMessage{
			Id:     int(message.Id),
			Date:   message.Date,
			Text:   message.Text,
			FromMe: message.FromMe,
		}

		response = append(response, msg)
	}

	return response, nil
}

func (s service) GetAllMessagesRestaurant(ctx context.Context, Uid, Rid int) ([]models.InfoMessage, error) {
	messages, err := s.chatService.GetAllMessagesRestaurant(ctx, &protoChat.Id{Uid: int32(Uid), Rid: int32(Rid)})
	if err != nil {
		return nil, err
	}

	response := make([]models.InfoMessage, 0)
	for _, message := range messages.More {
		msg := models.InfoMessage{
			Id:     int(message.Id),
			Date:   message.Date,
			Text:   message.Text,
			FromMe: message.FromMe,
		}

		response = append(response, msg)
	}

	return response, nil
}
