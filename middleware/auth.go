package middleware

import (
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
		ctx := models.GetContext(c)
		session, err := c.Cookie(config.SessionCookie)

		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		uid, ok, role := m.SessionUcase.Check(ctx, session.Value)

		if !ok {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if role == config.RoleUser { // тут проверяются права именно на обычного юзера
			user, err := m.UserUcase.GetByUid(ctx, uid)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			user.Uid = uid
			c.Set("User", user)
			return next(c)
		}

		return models.SendRedirectLogin(c)
	}
}

func InitMiddleware(userUcase userModel.UserUsecase, sessionUcase sessionModel.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase: sessionUcase,
		UserUcase:    userUcase,
	}
}
