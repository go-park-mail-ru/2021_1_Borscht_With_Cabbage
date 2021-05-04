package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/session"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	ChatUsecase    chat.ChatUsecase
	SessionUsecase session.SessionUsecase
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewChatHandler(chatUsecase chat.ChatUsecase,
	sessionUsecase session.SessionUsecase) chat.ChatHandler {
	return &chatHandler{
		ChatUsecase:    chatUsecase,
		SessionUsecase: sessionUsecase,
	}
}

func (ch chatHandler) GetKey(c echo.Context) error {
	ctx := models.GetContext(c)

	resuslt, err := ch.SessionUsecase.CreateKey(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	return models.SendResponse(c, &models.Key{Key: resuslt})
}

func (ch chatHandler) Connect(c echo.Context) error {
	ctx := models.GetContext(c)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, failError)
		return failError
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Open websocket": ws.RemoteAddr()})
	err = ch.ChatUsecase.Connect(ctx, ws)
	if err != nil {
		return err
	}

	defer func() {
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Close websocket": ws.RemoteAddr()})
		ch.ChatUsecase.UnConnect(ctx, ws)
	}()

	return ch.startRead(ctx, ws)
}

func (ch chatHandler) startRead(ctx context.Context, ws *websocket.Conn) error {
	for {
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"start read": ""})
		msg := new(models.FromClient)
		err := ws.ReadJSON(msg)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"error": failError.Error()})
			return failError
		}
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"new message": msg})
		err = ch.ChatUsecase.MessageCame(ctx, ws, *msg)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, failError)
			return failError
		}
	}
}

func (ch chatHandler) GetAllChats(c echo.Context) error {
	ctx := models.GetContext(c)

	result, err := ch.ChatUsecase.GetAllChats(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	response := make([]models.Response, 0)
	for i := range result {
		response = append(response, &result[i])
	}
	return models.SendMoreResponse(c, response...)
}

func (ch chatHandler) GetAllMessages(c echo.Context) error {
	ctx := models.GetContext(c)

	idStr := c.Param("id")
	if idStr == "" {
		sentErr := errors.BadRequestError("Error with id")
		logger.DeliveryLevel().ErrorLog(ctx, sentErr)
		return models.SendResponseWithError(c, sentErr)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sentErr := errors.BadRequestError("Error with id")
		logger.DeliveryLevel().ErrorLog(ctx, sentErr)
		return models.SendResponseWithError(c, sentErr)
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"id": id})

	result, err := ch.ChatUsecase.GetAllMessages(ctx, id)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	return models.SendMoreResponse(c, result)
}
