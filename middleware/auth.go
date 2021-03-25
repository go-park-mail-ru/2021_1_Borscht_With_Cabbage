package middleware

import (
	"fmt"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	SessionUcase sessionModel.SessionUsecase
	UserUcase    userModel.UserUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := c.Cookie(config.SessionCookie)

		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		phone, ok := m.SessionUcase.Check(c.Request().Context(), session.Value)
		if !ok {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		user, err := m.UserUcase.GetByNumber(c.Request().Context(), phone)
		c.Set("User", user)

		fmt.Println("THIS USER:", c.Get("User"))
		return next(c)
	}
}

func InitMiddleware(userUcase userModel.UserUsecase, sessionUcase sessionModel.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase: sessionUcase,
		UserUcase:    userUcase,
	}
}
