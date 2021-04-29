package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/services/auth"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type UserAuthMiddleware struct {
	AuthService auth.ServiceAuth
}

func (m *UserAuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)
		session, err := c.Cookie(config.SessionCookie)
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

		if sessionData.Role != config.RoleUser { // тут проверяются права именно на обычного юзера
			return models.SendRedirectLogin(c)
		}

		user, err := m.AuthService.GetByUid(ctx, sessionData.Id)
		if err != nil {
			return models.SendRedirectLogin(c)
		}
		user.Uid = sessionData.Id
		c.Set("User", user.User)
		logger.MiddleLevel().DataLog(ctx, "user auth", c.Get("User"))

		return next(c)
	}
}

func InitUserMiddleware(authService auth.ServiceAuth) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		AuthService: authService,
	}
}
