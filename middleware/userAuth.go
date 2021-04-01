package middleware

import (
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	sessionModel "github.com/borscht/backend/internal/session"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/labstack/echo/v4"
)

type UserAuthMiddleware struct {
	SessionUcase sessionModel.SessionUsecase
	UserUcase    userModel.UserUsecase
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
		*sessionData, exists, err = m.SessionUcase.Check(ctx, session.Value)
		if err != nil {
			return models.SendResponseWithError(c, err)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role != config.RoleUser { // тут проверяются права именно на обычного юзера
			return models.SendRedirectLogin(c)
		}

		user, err := m.UserUcase.GetByUid(ctx, sessionData.Id)
		if err != nil {
			return models.SendRedirectLogin(c)
		}
		user.Uid = sessionData.Id
		c.Set("User", user)
		return next(c)
	}
}

func InitUserMiddleware(userUcase userModel.UserUsecase, sessionUcase sessionModel.SessionUsecase) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		SessionUcase: sessionUcase,
		UserUcase:    userUcase,
	}
}
