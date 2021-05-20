package internal

import (
	"context"
	"database/sql"
	"testing"

	"github.com/borscht/backend/configProject"
	mocks2 "github.com/borscht/backend/services/chat/mocks"
	"github.com/borscht/backend/services/chat/models"
	proto "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
}

func TestService_GetAllChatsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	infoUser := proto.InfoUser{Id: 1, Role: configProject.RoleUser}
	chatsResult := []models.ChatInfo{{
		Message: models.Message{
			Mid:  1,
			Text: "hi",
			Date: "21.01.21",
		},
		User: models.User{
			Id:   1,
			Role: configProject.RoleUser,
		},
	}}

	user := models.User{Id: infoUser.Id, Role: configProject.RoleUser}
	chatRepoMock.EXPECT().GetAllChats(ctx, user).Return(chatsResult, nil)

	infoMessage, err := chatService.GetAllChats(ctx, &infoUser)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.More[0].Text, "hi")
}

func TestService_GetAllChatsRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	infoUser := proto.InfoUser{Id: 1, Role: configProject.RoleAdmin}
	chatsResult := []models.ChatInfo{{
		Message: models.Message{
			Mid:  1,
			Text: "hi",
			Date: "21.01.21",
		},
		User: models.User{
			Id:   1,
			Role: configProject.RoleAdmin,
		},
	}}

	user := models.User{Id: infoUser.Id, Role: configProject.RoleAdmin}
	chatRepoMock.EXPECT().GetAllChats(ctx, user).Return(chatsResult, nil)

	infoMessage, err := chatService.GetAllChats(ctx, &infoUser)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.More[0].Text, "hi")
}

func TestService_GetAllChats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	chatRepoMock.EXPECT().GetAllChats(ctx, models.User{}).Return([]models.ChatInfo{}, sql.ErrNoRows)

	id := proto.InfoUser{}
	_, err := chatService.GetAllChats(ctx, &id)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_GetAllMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	users := proto.Speakers{
		Speaker1: &proto.InfoUser{
			Id: 1, Role: configProject.RoleUser,
		},
		Speaker2: &proto.InfoUser{
			Id: 1, Role: configProject.RoleAdmin,
		},
	}
	chat := []models.Chat{{
		Sender: models.User{
			Id:   1,
			Role: configProject.RoleUser,
		},
		Recipient: models.User{
			Id:   1,
			Role: configProject.RoleAdmin,
		},
		Message: models.Message{
			Mid:  12,
			Text: "hi",
			Date: "01.01.21",
		},
	}}

	chatRepoMock.EXPECT().GetAllMessages(ctx, models.User{
		Id:   users.Speaker1.Id,
		Role: users.Speaker1.Role,
	}, models.User{
		Id:   users.Speaker2.Id,
		Role: users.Speaker2.Role,
	}).Return(chat, nil)

	infoMessage, err := chatService.GetAllMessages(ctx, &users)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.More[0].Text, "hi")
}

func TestService_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	info := proto.InfoMessage{
		Text: "hi",
		Date: "01.01.21",
		Participants: &proto.Participants{
			Sender: &proto.InfoUser{
				Id:   1,
				Role: configProject.RoleUser,
			},
			Recipient: &proto.InfoUser{
				Id:   1,
				Role: configProject.RoleUser,
			},
		},
	}

	chatRepoMock.EXPECT().SaveMessage(ctx, models.Chat{
		Message: models.Message{
			Text: info.Text,
			Date: info.Date,
		},
		Sender: models.User{
			Id:   info.Participants.Sender.Id,
			Role: info.Participants.Sender.Role,
		},
		Recipient: models.User{
			Id:   info.Participants.Recipient.Id,
			Role: info.Participants.Recipient.Role,
		},
	}).Return(int32(1), nil)

	infoMessage, err := chatService.ProcessMessage(ctx, &info)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.Text, "hi")
}
