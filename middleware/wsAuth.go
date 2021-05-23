package middleware

import (
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/services/auth"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type WsAuthMiddleware struct {
	AuthService auth.ServiceAuth
}

func (m *WsAuthMiddleware) WsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)

		key := c.Param(configProject.KeyParams)
		if key == "" {
			return models.SendResponseWithError(c, errors.BadRequestError("Error with key"))
		}
		logger.MiddleLevel().InfoLog(ctx, logger.Fields{"key": key})

		sessionData, exists, err := m.AuthService.CheckKey(ctx, key)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendResponseWithError(c, errors.BadRequestError("Key not valid"))
		}

		if sessionData.Role == configProject.RoleUser {
			user, err := m.AuthService.GetByUid(ctx, sessionData.Id)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}
			user.Uid = sessionData.Id
			logger.MiddleLevel().InfoLog(ctx, logger.Fields{"Websocket User": user.User})
			c.Set("User", user.User)
		}

		if sessionData.Role == configProject.RoleAdmin {
			restaurant, err := m.AuthService.GetByRid(ctx, sessionData.Id)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}
			restaurant.ID = sessionData.Id
			logger.MiddleLevel().InfoLog(ctx, logger.Fields{"Websocket Restaurant": restaurant.RestaurantInfo})
			c.Set("Restaurant", restaurant.RestaurantInfo)
		}

		return next(c)
	}
}

func InitWsMiddleware(authService auth.ServiceAuth) *WsAuthMiddleware {

	return &WsAuthMiddleware{
		AuthService: authService,
	}
}
