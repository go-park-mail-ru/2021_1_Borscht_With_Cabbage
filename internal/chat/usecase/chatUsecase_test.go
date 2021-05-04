package usecase

import (
	"context"
	"github.com/borscht/backend/internal/chat/mocks"
	"github.com/borscht/backend/internal/models"
	serviceMocks "github.com/borscht/backend/internal/services/mocks"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"testing"
)

func TestChatUsecase_ConnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_ConnectRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_UnConnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_UnConnectRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_UnConnect_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_MessageCameFromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	user := models.User{
		Uid:    1,
		Name:   "Daria",
		Avatar: "image.jpg",
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)
	logger.InitLogger()

	message := models.FromClient{
		Action: "",
		Payload: models.FromClientPayload{
			To: models.WsOpponent{
				Id:     1,
				Name:   "Daria",
				Avatar: "image.jpg",
			},
			Message: models.WsMessage{
				Id:   1,
				Date: "01.01.20 21:11",
				Text: "hi",
			},
		},
	}
	saveInfo := models.WsMessageForRepo{
		Date:       message.Payload.Message.Date,
		Content:    message.Payload.Message.Text,
		SentToId:   message.Payload.To.Id,
		SentFromId: 1,
	}

	ChatRepoMock.EXPECT().SaveMessageFromUser(ctx, saveInfo).Return(1, nil)

	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, message)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_MessageCameFromRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)
	logger.InitLogger()

	message := models.FromClient{
		Action: "",
		Payload: models.FromClientPayload{
			To: models.WsOpponent{
				Id:     1,
				Name:   "Daria",
				Avatar: "image.jpg",
			},
			Message: models.WsMessage{
				Id:   1,
				Date: "01.01.20 21:11",
				Text: "hi",
			},
		},
	}
	saveInfo := models.WsMessageForRepo{
		Date:       message.Payload.Message.Date,
		Content:    message.Payload.Message.Text,
		SentToId:   message.Payload.To.Id,
		SentFromId: 1,
	}

	ChatRepoMock.EXPECT().SaveMessageFromRestaurant(ctx, saveInfo).Return(1, nil)

	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, message)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_MessageCame_AuthError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)
	logger.InitLogger()

	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, models.FromClient{})
	if err == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_MessageFromUser(t *testing.T) {

}

func TestChatUsecase_MessageFromRestaurant(t *testing.T) {

}

func TestChatUsecase_GetAllChatsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)
	logger.InitLogger()

	ChatServiceMock.EXPECT().GetAllChats(ctx, 1, 0).Return([]models.BriefInfoChat{}, nil)

	_, err := chatUsecase.GetAllChats(ctx)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllChatsRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)
	logger.InitLogger()

	ChatServiceMock.EXPECT().GetAllChats(ctx, 0, 1).Return([]models.BriefInfoChat{}, nil)

	_, err := chatUsecase.GetAllChats(ctx)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllChats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)
	logger.InitLogger()

	_, err := chatUsecase.GetAllChats(ctx)
	if err == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllMessagesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	id := 1
	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)
	logger.InitLogger()

	ChatServiceMock.EXPECT().GetAllMessagesUser(ctx, user.Uid, id).Return([]models.InfoMessage{}, nil)
	ChatRepoMock.EXPECT().GetRestaurant(ctx, id).Return(&models.InfoOpponent{}, nil)

	_, err := chatUsecase.GetAllMessages(ctx, id)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllMessagesRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatRepoMock := mocks.NewMockChatRepo(ctrl)
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	chatUsecase := NewChatUsecase(ChatRepoMock, ChatServiceMock)

	id := 1
	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)
	logger.InitLogger()

	ChatServiceMock.EXPECT().GetAllMessagesRestaurant(ctx, restaurant.ID, id).Return([]models.InfoMessage{}, nil)
	ChatRepoMock.EXPECT().GetUser(ctx, id).Return(&models.InfoOpponent{}, nil)

	_, err := chatUsecase.GetAllMessages(ctx, id)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}
