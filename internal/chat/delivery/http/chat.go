package http

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/services/auth"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	ChatUsecase    chat.ChatUsecase
	SessionService auth.ServiceAuth
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			ctx := context.Background()
			origin := r.Header["Origin"]
			if len(origin) == 0 {
				err := errors.BadRequestError("host == 0")
				logger.DeliveryLevel().ErrorLog(ctx, err)
				return false
			}
			u, err := url.Parse(origin[0])
			if err != nil {
				logger.DeliveryLevel().ErrorLog(ctx, err)
				return false
			}
			client, err := url.Parse(config.Client)
			if err != nil {
				logger.DeliveryLevel().ErrorLog(ctx, err)
				return false
			}

			logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Host": u.Host, "ClientHost": client.Host})
			return u.Host == client.Host
		},
	}
)

func NewChatHandler(chatUsecase chat.ChatUsecase,
	sessionService auth.ServiceAuth) chat.ChatHandler {
	return &chatHandler{
		ChatUsecase:    chatUsecase,
		SessionService: sessionService,
	}
}

func getSessionInfo(ctx context.Context) (*models.SessionInfo, error) {
	userInterface := ctx.Value("User")
	if userInterface != nil {
		user, ok := userInterface.(models.User)
		if !ok {
			failError := errors.FailServerError("failed to convert to models.User")
			logger.UsecaseLevel().ErrorLog(ctx, failError)
			return nil, failError
		}

		return &models.SessionInfo{
			Id:   user.Uid,
			Role: configProject.RoleUser,
		}, nil
	}

	restaurantInterface := ctx.Value("Restaurant")
	if restaurantInterface != nil {
		restaurant, ok := restaurantInterface.(models.RestaurantInfo)
		if !ok {
			failError := errors.FailServerError("failed to convert to models.Restaurant")
			logger.UsecaseLevel().ErrorLog(ctx, failError)
			return nil, failError
		}
		return &models.SessionInfo{
			Id:   restaurant.ID,
			Role: configProject.RoleAdmin,
		}, nil
	}

	return nil, errors.BadRequestError("not authorization")
}

func (ch chatHandler) GetKey(c echo.Context) error {
	ctx := models.GetContext(c)

	sessionInfo, err := getSessionInfo(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}

	resuslt, err := ch.SessionService.CreateKey(ctx, *sessionInfo)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	return models.SendResponse(c, &models.Key{Key: resuslt})
}

func (ch chatHandler) Connect(c echo.Context) error {
	ctx := models.GetContext(c)

	logger.DeliveryLevel().InlineInfoLog(ctx, "!!!!!!!!!!!")
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
		// read message
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"start read": ""})
		msg := new(models.FromClient)
		err := ws.ReadJSON(msg)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"error": failError.Error()})
			return failError
		}

		// message processing
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"new message": msg})
		err = ch.ChatUsecase.ProcessMessage(ctx, ws, *msg)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, failError)
			return failError
		}
	}

	errFail := errors.FailServerError("breake read message in websocket")
	logger.DeliveryLevel().ErrorLog(ctx, errFail)
	return errFail
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
