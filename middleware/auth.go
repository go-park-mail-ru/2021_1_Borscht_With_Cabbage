package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	SessionUcase         sessionModel.SessionUsecase
	UserUcase            userModel.UserUsecase
	RestaurantAdminUcase restaurantAdmin.AdminUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := c.Cookie(config.SessionCookie)
		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		sessionData := new(models.SessionInfo)
		var exists bool
		*sessionData, exists, err = m.SessionUcase.Check(session.Value)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role == config.RoleUser {
			user, err := m.UserUcase.GetByUid(sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			user.Uid = sessionData.Id
			c.Set("User", user)
		}

		if sessionData.Role == config.RoleAdmin {
			restaurant, err := m.RestaurantAdminUcase.GetByRid(sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			restaurant.ID = sessionData.Id
			c.Set("Restaurant", restaurant)
		}

		return next(c)
	}
}

func InitAuthMiddleware(userUcase userModel.UserUsecase, restaurantAdminUcase restaurantAdmin.AdminUsecase, sessionUcase sessionModel.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase:         sessionUcase,
		UserUcase:            userUcase,
		RestaurantAdminUcase: restaurantAdminUcase,
	}
}
