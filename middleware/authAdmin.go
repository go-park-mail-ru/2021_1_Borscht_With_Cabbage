package middleware

import (
	"fmt"
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
		session, err := c.Cookie(config.SessionCookie)

		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		id, ok, role := m.SessionUcase.Check(session.Value)

		if !ok {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if role == config.RoleAdmin { // тут проверяются права именно на обычного юзера
			restaurant, err := m.AdminUcase.GetByRid(id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			restaurant.ID = id
			c.Set("Restaurant", restaurant)
			fmt.Println("THIS RESTAURANT:", c.Get("Restaurant"))
			return next(c)
		}

		return models.SendRedirectLogin(c)
	}
}

func InitAdminMiddleware(adminUcase adminModel.AdminUsecase, sessionUcase sessionModel.SessionUsecase) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		SessionUcase: sessionUcase,
		AdminUcase:   adminUcase,
	}
}
