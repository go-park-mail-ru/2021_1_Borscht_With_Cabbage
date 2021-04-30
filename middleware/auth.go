package middleware

import (
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/services/auth"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	AuthService auth.ServiceAuth
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
		*sessionData, exists, err = m.AuthService.CheckSession(ctx, session.Value)
		fmt.Println("session data: ", sessionData)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role == config.RoleUser {
			user, err := m.AuthService.GetByUid(ctx, sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			user.Uid = sessionData.Id
			c.Set("User", user.User)
		}

		if sessionData.Role == config.RoleAdmin {
			restaurant, err := m.AuthService.GetByRid(ctx, sessionData.Id)
			fmt.Println("rest ", restaurant)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			restaurant.ID = sessionData.Id
			c.Set("Restaurant", restaurant.RestaurantInfo)
		}

		return next(c)
	}
}

func InitAuthMiddleware(authService auth.ServiceAuth) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}
