package middleware

import (
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"

	//adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/internal/services/auth"
	//sessionModel "github.com/borscht/backend/internal/session"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type AdminAuthMiddleware struct {
	AuthService auth.ServiceAuth
}

func (m *AdminAuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)
		session, err := c.Cookie(configProject.SessionCookie)
		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		sessionData := new(models.SessionInfo)
		var exists bool
		*sessionData, exists, err = m.AuthService.CheckSession(ctx, session.Value)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role != configProject.RoleAdmin { // тут проверяются права именно на обычного юзера
			return models.SendRedirectLogin(c)
		}

		restaurant, err := m.AuthService.GetByRid(ctx, sessionData.Id)
		if err != nil {
			return models.SendRedirectLogin(c)
		}
		restaurant.ID = sessionData.Id
		c.Set("Restaurant", restaurant.RestaurantInfo)
		logger.MiddleLevel().DataLog(ctx, "restaurant auth", c.Get("Restaurant"))

		return next(c)
	}
}

func InitAdminMiddleware(authService auth.ServiceAuth) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		AuthService: authService,
	}
}
