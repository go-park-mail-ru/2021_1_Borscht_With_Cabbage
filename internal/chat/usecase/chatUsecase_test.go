package usecase

import (
	"context"
	"testing"

	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	serviceMocks "github.com/borscht/backend/internal/services/mocks"
	"github.com/borscht/backend/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
)

func TestChatUsecase_ConnectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", 1)

	err := chatUsecase.Connect(ctx, &websocket.Conn{})
	if err == nil {
		t.Errorf("incorrect result")
		return
	}
}

//
//func TestChatUsecase_MessageCameFromUser(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
//	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
//	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)
//
//	user := models.User{
//		Uid:    1,
//		Name:   "Daria",
//		Avatar: "image.jpg",
//	}
//	c := context.Background()
//	ctx := context.WithValue(c, "User", user)
//	logger.InitLogger()
//
//	msg := models.FromClient{
//		Action: "",
//		Payload: models.FromClientPayload{
//			To: models.WsOpponent{
//				Id:     1,
//				Name:   "Daria",
//				Avatar: "image.jpg",
//			},
//			Message: models.WsMessage{
//				Id:   1,
//				Date: "01.01.20 21:11",
//				Text: "hi",
//			},
//		},
//	}
//	message := models.ChatMessage{
//		Text: msg.Payload.Message.Text,
//		Date: msg.Payload.Message.Date,
//	}
//	me := models.ChatUser{
//		Id:   user.Uid,
//		Role: config.RoleUser,
//	}
//	opponent := models.ChatUser{
//		Id:   1,
//		Role: config.RoleAdmin,
//	}
//
//	infoChat := models.InfoChatMessage{
//		Message:   message,
//		Sender:    me,
//		Recipient: opponent,
//	}
//
//	ChatServiceMock.EXPECT().ProcessMessage(ctx, infoChat).Return(models.InfoChatMessage{}, nil)
//
//	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, msg)
//	if err != nil {
//		t.Errorf("incorrect result")
//		return
//	}
//}
//
//func TestChatUsecase_MessageCameFromRestaurant(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
//	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
//	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)
//
//	restaurant := models.RestaurantInfo{
//		ID: 1,
//	}
//	c := context.Background()
//	ctx := context.WithValue(c, "Restaurant", restaurant)
//	logger.InitLogger()
//
//	msg := models.FromClient{
//		Action: "",
//		Payload: models.FromClientPayload{
//			To: models.WsOpponent{
//				Id:     1,
//				Name:   "Daria",
//				Avatar: "image.jpg",
//			},
//			Message: models.WsMessage{
//				Id:   1,
//				Date: "01.01.20 21:11",
//				Text: "hi",
//			},
//		},
//	}
//	message := models.ChatMessage{
//		Text: msg.Payload.Message.Text,
//		Date: msg.Payload.Message.Date,
//	}
//	me := models.ChatUser{
//		Id:   restaurant.ID,
//		Role: config.RoleAdmin,
//	}
//	opponent := models.ChatUser{
//		Id:   1,
//		Role: config.RoleUser,
//	}
//	infoChat := models.InfoChatMessage{
//		Message:   message,
//		Sender:    me,
//		Recipient: opponent,
//	}
//
//	ChatServiceMock.EXPECT().ProcessMessage(ctx, infoChat).Return(models.InfoChatMessage{}, nil)
//
//	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, msg)
//	if err != nil {
//		t.Errorf("incorrect result")
//		return
//	}
//}
//
//func TestChatUsecase_MessageCame_AuthError(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
//	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
//	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)
//
//	c := context.Background()
//	ctx := context.WithValue(c, "Restaurant", 1)
//	logger.InitLogger()
//
//	err := chatUsecase.MessageCame(ctx, &websocket.Conn{}, models.FromClient{})
//	if err == nil {
//		t.Errorf("incorrect result")
//		return
//	}
//}

func TestChatUsecase_GetAllChatsUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)
	logger.InitLogger()

	me := models.ChatUser{
		Id:   user.Uid,
		Role: configProject.RoleUser,
	}
	ChatServiceMock.EXPECT().GetAllChats(ctx, me).Return([]models.InfoChatMessage{}, nil)

	_, err := chatUsecase.GetAllChats(ctx)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllChatsRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)
	logger.InitLogger()

	me := models.ChatUser{
		Id:   restaurant.ID,
		Role: configProject.RoleAdmin,
	}
	ChatServiceMock.EXPECT().GetAllChats(ctx, me).Return([]models.InfoChatMessage{}, nil)

	_, err := chatUsecase.GetAllChats(ctx)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllChats_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

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
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

	id := 1
	user := models.User{
		Uid: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "User", user)
	logger.InitLogger()

	me := models.ChatUser{
		Id:   user.Uid,
		Role: configProject.RoleUser,
	}
	opponent := models.ChatUser{
		Id:   id,
		Role: configProject.RoleAdmin,
	}

	AuthServiceMock.EXPECT().GetByRid(ctx, id).Return(&models.SuccessRestaurantResponse{}, nil)
	ChatServiceMock.EXPECT().GetAllMessages(ctx, me, opponent).Return([]models.InfoChatMessage{}, nil)

	_, err := chatUsecase.GetAllMessages(ctx, id)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestChatUsecase_GetAllMessagesRestaurant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ChatServiceMock := serviceMocks.NewMockServiceChat(ctrl)
	AuthServiceMock := serviceMocks.NewMockServiceAuth(ctrl)
	chatUsecase := NewChatUsecase(ChatServiceMock, AuthServiceMock)

	id := 1
	restaurant := models.RestaurantInfo{
		ID: 1,
	}
	c := context.Background()
	ctx := context.WithValue(c, "Restaurant", restaurant)
	logger.InitLogger()
	me := models.ChatUser{
		Id:   restaurant.ID,
		Role: configProject.RoleAdmin,
	}
	opponent := models.ChatUser{
		Id:   id,
		Role: configProject.RoleUser,
	}

	AuthServiceMock.EXPECT().GetByUid(ctx, id).Return(&models.SuccessUserResponse{}, nil)
	ChatServiceMock.EXPECT().GetAllMessages(ctx, me, opponent).Return([]models.InfoChatMessage{}, nil)

	_, err := chatUsecase.GetAllMessages(ctx, id)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}
