package middleware

import (
	"backend/api/domain"
	"backend/api/domain/user"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	SessionUcase domain.SessionUsecase
	UserUcase    user.UserUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*domain.CustomContext)

		session, err := c.Cookie(domain.SessionCookie)

		if err != nil {
			cc.User = nil
			return next(c) // пользователь не вошел
		}

		phone, ok := m.SessionUcase.Check(session.Value, cc)
		if !ok {
			cc.User = nil
			return next(c) // пользователь не вошел
		}

		user, err := m.UserUcase.GetByNumber(cc, phone)
		cc.User = &user

		return next(c)
	}
}

func InitMiddleware(uus user.UserUsecase, sus domain.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase: sus,
		UserUcase: uus,
	}
}
