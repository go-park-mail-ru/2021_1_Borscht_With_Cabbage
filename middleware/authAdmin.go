package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	sessionModel "github.com/borscht/backend/internal/session"
	"github.com/labstack/echo/v4"
)

type AdminAuthMiddleware struct {
	SessionUcase sessionModel.SessionUsecase
	AdminUcase   adminModel.AdminUsecase
}

func (m *AdminAuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)
		session, err := c.Cookie(config.SessionCookie)
		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		sessionData := new(models.SessionInfo)
		var exists bool
		*sessionData, exists, err = m.SessionUcase.Check(ctx, session.Value)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role != config.RoleAdmin { // тут проверяются права именно на обычного юзера
			return models.SendRedirectLogin(c)
		}

		restaurant, err := m.AdminUcase.GetByRid(ctx, sessionData.Id)
		if err != nil {
			return models.SendRedirectLogin(c)
		}
		restaurant.ID = sessionData.Id
		c.Set("Restaurant", restaurant)
		return next(c)
	}
}

func InitAdminMiddleware(adminUcase adminModel.AdminUsecase, sessionUcase sessionModel.SessionUsecase) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		SessionUcase: sessionUcase,
		AdminUcase:   adminUcase,
	}
}
