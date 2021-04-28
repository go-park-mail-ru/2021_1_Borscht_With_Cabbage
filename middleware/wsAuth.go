package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type WsAuthMiddleware struct {
	SessionUcase         sessionModel.SessionUsecase
	UserUcase            userModel.UserUsecase
	RestaurantAdminUcase restaurantAdmin.AdminRestaurantUsecase
}

func (m *WsAuthMiddleware) WsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)

		key := c.Param(config.KeyParams)
		if key == "" {
			return models.SendResponseWithError(c, errors.BadRequestError("Error with key"))
		}
		logger.MiddleLevel().InfoLog(ctx, logger.Fields{"key": key})

		sessionData, exists, err := m.SessionUcase.CheckKey(ctx, key)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendResponseWithError(c, errors.BadRequestError("Key not valid"))
		}

		if sessionData.Role == config.RoleUser {
			user, err := m.UserUcase.GetByUid(ctx, sessionData.Id)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}
			user.Uid = sessionData.Id
			logger.MiddleLevel().InfoLog(ctx, logger.Fields{"Websocket User": user.User})
			c.Set("User", user.User)
		}

		if sessionData.Role == config.RoleAdmin {
			restaurant, err := m.RestaurantAdminUcase.GetByRid(ctx, sessionData.Id)
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

func InitWsMiddleware(userUcase userModel.UserUsecase,
	restaurantAdminUcase restaurantAdmin.AdminRestaurantUsecase,
	sessionUcase sessionModel.SessionUsecase) *WsAuthMiddleware {

	return &WsAuthMiddleware{
		SessionUcase:         sessionUcase,
		UserUcase:            userUcase,
		RestaurantAdminUcase: restaurantAdminUcase,
	}
}
