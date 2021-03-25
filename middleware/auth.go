package middleware

import (
	"fmt"
	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	_sessionModel "github.com/borscht/backend/internal/session"
	_userModel "github.com/borscht/backend/internal/user"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	SessionUcase _sessionModel.SessionUsecase
	UserUcase    _userModel.UserUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*models.CustomContext)

		session, err := c.Cookie(config.SessionCookie)

		if err != nil {
			cc.User = nil
			return next(c) // пользователь не вошел
		}

		uid, ok := m.SessionUcase.Check(session.Value)
		if !ok {
			cc.User = nil
			return next(c) // пользователь не вошел
		}

		user, err := m.UserUcase.GetByUid(uid)
		cc.User = &user
		cc.User.Uid = uid

		fmt.Println("THIS USER:", user)
		return next(cc)
	}
}

func InitMiddleware(userUcase _userModel.UserUsecase, sessionUcase _sessionModel.SessionUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUcase: sessionUcase,
		UserUcase:    userUcase,
	}
}
