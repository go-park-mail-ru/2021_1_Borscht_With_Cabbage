package http

import (
	"strconv"

	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	ChatUsecase chat.ChatUsecase
}

func NewChatHandler(chatUsecase chat.ChatUsecase) chat.ChatHandler {
	return &chatHandler{
		ChatUsecase: chatUsecase,
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
