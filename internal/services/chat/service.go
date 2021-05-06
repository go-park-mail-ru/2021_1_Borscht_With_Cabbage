package chat

import (
	"context"

	"github.com/borscht/backend/internal/models"
	protoChat "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/logger"
)

type ServiceChat interface {
	ProcessMessage(ctx context.Context, chat models.InfoChatMessage) (models.InfoChatMessage, error)
	GetAllChats(ctx context.Context, user models.ChatUser) ([]models.InfoChatMessage, error)
	GetAllMessages(ctx context.Context, speaker1, speaker2 models.ChatUser) ([]models.InfoChatMessage, error)
}

type service struct {
	chatService protoChat.ChatClient
}

func NewService(chatService protoChat.ChatClient) ServiceChat {
	return &service{
		chatService: chatService,
	}
}

func (s service) GetAllChats(ctx context.Context, user models.ChatUser) (
	[]models.InfoChatMessage, error) {

	chats, err := s.chatService.GetAllChats(ctx, &protoChat.InfoUser{Id: int32(user.Id), Role: user.Role})
	if err != nil {
		return nil, err
	}

	response := make([]models.InfoChatMessage, 0)
	for _, chat := range chats.More {

		chatModel := convertMessage(chat)
		response = append(response, chatModel)
	}

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"answer from microservice chat": response})
	return response, nil
}

func (s service) GetAllMessages(ctx context.Context, speaker1, speaker2 models.ChatUser) (
	[]models.InfoChatMessage, error) {

	protoSpeaker1 := protoChat.InfoUser{Id: int32(speaker1.Id), Role: speaker1.Role}
	protoSpeaker2 := protoChat.InfoUser{Id: int32(speaker2.Id), Role: speaker2.Role}
	speakers := protoChat.Speakers{Speaker1: &protoSpeaker1, Speaker2: &protoSpeaker2}
	messages, err := s.chatService.GetAllMessages(ctx, &speakers)
	if err != nil {
		return nil, err
	}

	response := make([]models.InfoChatMessage, 0)
	for _, message := range messages.More {

		msg := convertMessage(message)
		response = append(response, msg)
	}

	return response, nil
}

func convertMessage(message *protoChat.InfoMessage) models.InfoChatMessage {
	sender := models.ChatUser{
		Id:   int(message.Participants.Sender.Id),
		Role: message.Participants.Sender.Role,
	}
	recipient := models.ChatUser{
		Id:   int(message.Participants.Recipient.Id),
		Role: message.Participants.Recipient.Role,
	}
	msg := models.ChatMessage{
		Mid:  int(message.Id),
		Date: message.Date,
		Text: message.Text,
	}

	response := models.InfoChatMessage{
		Message:   msg,
		Sender:    sender,
		Recipient: recipient,
	}

	return response
}

func (s service) ProcessMessage(ctx context.Context, chat models.InfoChatMessage) (
	models.InfoChatMessage, error) {

	logger.UsecaseLevel().DebugLog(ctx, logger.Fields{"chat": chat})

	sender := protoChat.InfoUser{Id: int32(chat.Sender.Id), Role: chat.Sender.Role}
	recipient := protoChat.InfoUser{Id: int32(chat.Recipient.Id), Role: chat.Recipient.Role}
	participants := protoChat.Participants{Sender: &sender, Recipient: &recipient}
	request := protoChat.InfoMessage{
		Id:           int32(chat.Message.Mid),
		Date:         chat.Message.Date,
		Text:         chat.Message.Text,
		Participants: &participants,
	}

	message, err := s.chatService.ProcessMessage(ctx, &request)
	if err != nil {
		return models.InfoChatMessage{}, err
	}

	response := convertMessage(message)
	return response, nil
}
