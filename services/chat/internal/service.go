package internal

import (
	"context"
	"sort"

	"github.com/borscht/backend/services/chat/models"
	chatRepo "github.com/borscht/backend/services/chat/repository"
	protoChat "github.com/borscht/backend/services/proto/chat"
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

func (s service) GetAllChats(ctx context.Context, info *protoChat.InfoUser) (*protoChat.MoreInfoMessage, error) {
	logger.UsecaseLevel().InfoLog(ctx, logger.Fields{"user": info})

	result := new(protoChat.MoreInfoMessage)
	fromMe, err := s.chatRepo.GetAllChatsFromUser(ctx, models.User{
		Id:   info.Id,
		Role: info.Role,
	})
	if err != nil {
		return nil, err
	}

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"fromMe": fromMe})
	for _, value := range fromMe {
		recipient := protoChat.InfoUser{
			Id:   value.User.Id,
			Role: value.User.Role,
		}
		proto := protoChat.InfoMessage{
			Id:   value.Message.Mid,
			Text: value.Message.Text,
			Date: value.Message.Date,
		}
		pr := protoChat.Participants{
			Sender:    info,
			Recipient: &recipient,
		}
		proto.Participants = &pr
		result.More = append(result.More, &proto)
	}

	toMe, err := s.chatRepo.GetAllChatsToUser(ctx, models.User{
		Id:   info.Id,
		Role: info.Role,
	})
	if err != nil {
		return nil, err
	}

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"toMe": toMe})
	for _, value := range toMe {
		sender := protoChat.InfoUser{
			Id:   value.User.Id,
			Role: value.User.Role,
		}
		proto := protoChat.InfoMessage{
			Id:   value.Message.Mid,
			Text: value.Message.Text,
			Date: value.Message.Date,
		}
		pr := protoChat.Participants{
			Sender:    &sender,
			Recipient: info,
		}
		proto.Participants = &pr
		result.More = append(result.More, &proto)
	}

	sort.SliceStable(result.More, func(i, j int) bool {
		return result.More[i].Id > result.More[j].Id
	})

	logger.UsecaseLevel().InfoLog(ctx, logger.Fields{"result": result})
	return result, nil
}

func (ch service) GetAllMessages(ctx context.Context, users *protoChat.Speakers) (
	result *protoChat.MoreInfoMessage, err error) {

	chat, err := ch.chatRepo.GetAllMessages(ctx, models.User{
		Id:   users.Speaker1.Id,
		Role: users.Speaker1.Role,
	}, models.User{
		Id:   users.Speaker2.Id,
		Role: users.Speaker2.Role,
	})

	if err != nil {
		return nil, err
	}

	result = new(protoChat.MoreInfoMessage)
	for _, value := range chat {
		sender := protoChat.InfoUser{
			Id:   value.Sender.Id,
			Role: value.Sender.Role,
		}
		recipient := protoChat.InfoUser{
			Id:   value.Recipient.Id,
			Role: value.Recipient.Role,
		}
		pr := protoChat.Participants{
			Sender:    &sender,
			Recipient: &recipient,
		}
		proto := protoChat.InfoMessage{
			Id:           value.Message.Mid,
			Text:         value.Message.Text,
			Date:         value.Message.Date,
			Participants: &pr,
		}
		result.More = append(result.More, &proto)
	}

	return result, nil
}

func (ch service) SendMessage(ctx context.Context, info *protoChat.InfoMessage) (*protoChat.InfoMessage, error) {
	message := models.Message{
		Text: info.Text,
		Date: info.Date,
	}
	sender := models.User{
		Id:   info.Participants.Sender.Id,
		Role: info.Participants.Sender.Role,
	}
	recipient := models.User{
		Id:   info.Participants.Recipient.Id,
		Role: info.Participants.Recipient.Role,
	}
	mid, err := ch.chatRepo.SaveMessage(ctx, models.Chat{
		Message:   message,
		Sender:    sender,
		Recipient: recipient,
	})
	if err != nil {
		return nil, err
	}

	info.Id = mid
	return info, nil
}
