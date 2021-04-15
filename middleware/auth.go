package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	SessionUcase         sessionModel.SessionUsecase
	UserUcase            userModel.UserUsecase
	RestaurantAdminUcase restaurantAdmin.AdminRestaurantUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)
		logger.MiddleLevel().InlineDebugLog(ctx, "Autorization")
		session, err := c.Cookie(config.SessionCookie)
		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		logger.MiddleLevel().InlineDebugLog(ctx, session.Value)

		sessionData := new(models.SessionInfo)
		var exists bool
		*sessionData, exists, err = m.SessionUcase.Check(ctx, session.Value)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role == config.RoleUser {
			user, err := m.UserUcase.GetByUid(ctx, sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			user.Uid = sessionData.Id
			c.Set("User", user.User)
		}

		if sessionData.Role == config.RoleAdmin {
			restaurant, err := m.RestaurantAdminUcase.GetByRid(ctx, sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			restaurant.ID = sessionData.Id
			c.Set("Restaurant", restaurant.RestaurantInfo)
		}

		return next(c)
	}
}

func InitAuthMiddleware(userUcase userModel.UserUsecase, restaurantAdminUcase restaurantAdmin.AdminRestaurantUsecase, sessionUcase sessionModel.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase:         sessionUcase,
		UserUcase:            userUcase,
		RestaurantAdminUcase: restaurantAdminUcase,
	}
}
