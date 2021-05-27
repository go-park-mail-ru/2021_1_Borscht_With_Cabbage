package chat

import (
	"context"
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/services/chat/mocks"
	protoChat "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_GetAllChats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockChatClient(ctrl)
	chatService := NewService(clientMock)

	c := context.Background()
	logger.InitLogger()

	user := models.ChatUser{
		Id:   1,
		Role: configProject.RoleUser,
	}
	chat := []*protoChat.InfoMessage{{
		Id: 1,
		Participants: &protoChat.Participants{
			Sender: &protoChat.InfoUser{
				Id: 1,
			},
			Recipient: &protoChat.InfoUser{
				Id: 1,
			},
		},
	}}
	clientMock.EXPECT().GetAllChats(c, &protoChat.InfoUser{Id: int32(user.Id), Role: user.Role}).Return(&protoChat.MoreInfoMessage{
		More: chat,
	}, nil)

	_, err := chatService.GetAllChats(c, user)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_GetAllMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockChatClient(ctrl)
	chatService := NewService(clientMock)

	c := context.Background()
	logger.InitLogger()

	speaker1 := models.ChatUser{Id: 1, Role: configProject.RoleUser}
	speaker2 := models.ChatUser{Id: 1, Role: configProject.RoleUser}

	messages := protoChat.MoreInfoMessage{
		More: []*protoChat.InfoMessage{{
			Id: 1,
			Participants: &protoChat.Participants{
				Sender: &protoChat.InfoUser{
					Id: 1,
				},
				Recipient: &protoChat.InfoUser{
					Id: 1,
				},
			},
		},
		},
	}

	protoSpeaker1 := protoChat.InfoUser{Id: int32(speaker1.Id), Role: speaker1.Role}
	protoSpeaker2 := protoChat.InfoUser{Id: int32(speaker2.Id), Role: speaker2.Role}
	speakers := protoChat.Speakers{Speaker1: &protoSpeaker1, Speaker2: &protoSpeaker2}
	clientMock.EXPECT().GetAllMessages(c, &speakers).Return(&messages, nil)

	_, err := chatService.GetAllMessages(c, speaker1, speaker2)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mocks.NewMockChatClient(ctrl)
	chatService := NewService(clientMock)

	c := context.Background()
	logger.InitLogger()

	chat := models.InfoChatMessage{
		Message: models.ChatMessage{
			Mid:  1,
			Text: "hi",
		},
		Sender: models.ChatUser{
			Id: 1,
		},
		Recipient: models.ChatUser{
			Id: 1,
		},
	}
	sender := protoChat.InfoUser{Id: int32(chat.Sender.Id), Role: chat.Sender.Role}
	recipient := protoChat.InfoUser{Id: int32(chat.Recipient.Id), Role: chat.Recipient.Role}
	participants := protoChat.Participants{Sender: &sender, Recipient: &recipient}
	request := protoChat.InfoMessage{
		Id:           int32(chat.Message.Mid),
		Date:         chat.Message.Date,
		Text:         chat.Message.Text,
		Participants: &participants,
	}
	response := protoChat.InfoMessage{
		Id: 1,
		Participants: &protoChat.Participants{
			Sender: &protoChat.InfoUser{
				Id: 1,
			},
			Recipient: &protoChat.InfoUser{
				Id: 1,
			},
		},
	}

	clientMock.EXPECT().ProcessMessage(c, &request).Return(&response, nil)
	_, err := chatService.ProcessMessage(c, chat)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}
