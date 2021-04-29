package http

import (
	"context"
	"net/http"

	custWebsocket "github.com/borscht/backend/internal/chat"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/session"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	WsUsecase      custWebsocket.WebSocketUsecase
	SessionUsecase session.SessionUsecase
}

func NewWebSocketHandler(wsUsecase custWebsocket.WebSocketUsecase,
	sessionUsecase session.SessionUsecase) custWebsocket.WebSocketHandler {
	return &WebSocketHandler{
		WsUsecase:      wsUsecase,
		SessionUsecase: sessionUsecase,
	}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (w WebSocketHandler) GetKey(c echo.Context) error {
	ctx := models.GetContext(c)

	resuslt, err := w.SessionUsecase.CreateKey(ctx)
	if err != nil {
		return models.SendResponseWithError(c, err)
	}
	return models.SendResponse(c, &models.Key{Key: resuslt})
}

func (w WebSocketHandler) Connect(c echo.Context) error {
	ctx := models.GetContext(c)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.DeliveryLevel().ErrorLog(ctx, failError)
		return failError
	}
	logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Open websocket": ws.RemoteAddr()})
	err = w.WsUsecase.Connect(ctx, ws)
	if err != nil {
		return err
	}

	defer func() {
		logger.DeliveryLevel().InfoLog(ctx, logger.Fields{"Close websocket": ws.RemoteAddr()})
		w.WsUsecase.UnConnect(ctx, ws)
	}()

	return w.startRead(ctx, ws)
}

func (w WebSocketHandler) startRead(ctx context.Context, ws *websocket.Conn) error {
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
		err = w.WsUsecase.MessageCame(ctx, ws, *msg)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.DeliveryLevel().ErrorLog(ctx, failError)
			return failError
		}
	}
}
