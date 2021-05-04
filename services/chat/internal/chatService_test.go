package internal

import (
	"context"
	mocks2 "github.com/borscht/backend/services/chat/mocks"
	proto "github.com/borscht/backend/services/proto/chat"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
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

	id := proto.Id{Uid: 1}
	chatsResult := []*proto.BriefInfoChat{{LastMessage: "hi"}}

	chatRepoMock.EXPECT().GetAllChatsUser(ctx, 1).Return(chatsResult, nil)

	infoUser, err := chatService.GetAllChats(ctx, &id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoUser.More[0].LastMessage, "hi")
}

func TestService_GetAllChatsRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	id := proto.Id{Rid: 1}
	chatsResult := []*proto.BriefInfoChat{{LastMessage: "hi"}}

	chatRepoMock.EXPECT().GetAllChatsRestaurant(ctx, 1).Return(chatsResult, nil)

	infoUser, err := chatService.GetAllChats(ctx, &id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoUser.More[0].LastMessage, "hi")
}

func TestService_GetAllChats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	id := proto.Id{}

	_, err := chatService.GetAllChats(ctx, &id)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
}

func TestService_GetAllMessagesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	id := proto.Id{Uid: 1, Rid: 1}
	fromMe := []*proto.InfoMessage{{Text: "hi", Id: 1}}
	toMe := []*proto.InfoMessage{{Text: "hello", Id: 2}}

	chatRepoMock.EXPECT().GetAllMessagesFromUser(ctx, 1, 1).Return(fromMe, nil)
	chatRepoMock.EXPECT().GetAllMessagesFromRestaurant(ctx, 1, 1).Return(toMe, nil)

	infoMessage, err := chatService.GetAllMessagesUser(ctx, &id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.More[0].Text, "hello")
}

func TestService_GetAllMessagesRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chatRepoMock := mocks2.NewMockChatRepo(ctrl)
	chatService := NewService(chatRepoMock)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()

	id := proto.Id{Uid: 1, Rid: 1}
	fromMe := []*proto.InfoMessage{{Text: "hi", Id: 1}}
	toMe := []*proto.InfoMessage{{Text: "hello", Id: 2}}

	chatRepoMock.EXPECT().GetAllMessagesFromUser(ctx, 1, 1).Return(fromMe, nil)
	chatRepoMock.EXPECT().GetAllMessagesFromRestaurant(ctx, 1, 1).Return(toMe, nil)

	infoMessage, err := chatService.GetAllMessagesRestaurant(ctx, &id)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	require.EqualValues(t, infoMessage.More[0].Text, "hello")
}
