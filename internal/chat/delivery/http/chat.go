package http

import (
	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
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
