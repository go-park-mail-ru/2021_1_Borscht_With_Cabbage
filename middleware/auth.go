package middleware

import (
	"github.com/borscht/backend/configProject"
	"github.com/borscht/backend/internal/models"
	adminModel "github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/internal/services/auth"
	userModel "github.com/borscht/backend/internal/user"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	AuthService  auth.ServiceAuth
	UserUsecase  userModel.UserUsecase
	AdminUsecase adminModel.AdminRestaurantUsecase
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := models.GetContext(c)
		logger.MiddleLevel().InlineDebugLog(ctx, "Autorization")
		session, err := c.Cookie(configProject.SessionCookie)
		if err != nil {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		logger.MiddleLevel().InlineDebugLog(ctx, session.Value)

		sessionData := new(models.SessionInfo)
		var exists bool
		*sessionData, exists, err = m.AuthService.CheckSession(ctx, session.Value)
		logger.MiddleLevel().DebugLog(ctx, logger.Fields{"session data": sessionData})
		if err != nil {
			failErr := errors.FailServerError("CheckSession: " + err.Error())
			logger.MiddleLevel().ErrorLog(ctx, failErr)
			return models.SendResponseWithError(c, failErr)
		}
		if !exists {
			return models.SendRedirectLogin(c) // пользователь не вошел
		}

		if sessionData.Role == configProject.RoleUser {
			user, err := m.AuthService.GetByUid(ctx, sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			user.Uid = sessionData.Id

			address, err := m.UserUsecase.GetAddress(ctx, user.Uid)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}

			if address != nil {
				user.Address = *address
			}

			c.Set("User", user.User)
		}

		if sessionData.Role == configProject.RoleAdmin {
			restaurant, err := m.AuthService.GetByRid(ctx, sessionData.Id)
			if err != nil {
				return models.SendRedirectLogin(c)
			}
			restaurant.ID = sessionData.Id

			address, err := m.AdminUsecase.GetAddress(ctx, restaurant.ID)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}

			if address != nil {
				restaurant.Address = *address
			}

			categories, err := m.AdminUsecase.GetCategories(ctx, restaurant.ID)
			if err != nil {
				return models.SendResponseWithError(c, err)
			}
			restaurant.Categories = categories.CategoriesID

			c.Set("Restaurant", restaurant.RestaurantInfo)
		}

		return next(c)
	}
}

func InitAuthMiddleware(authService auth.ServiceAuth, userUsecase userModel.UserUsecase,
	adminUsecase adminModel.AdminRestaurantUsecase) *AuthMiddleware {

	return &AuthMiddleware{
		AuthService:  authService,
		UserUsecase:  userUsecase,
		AdminUsecase: adminUsecase,
	}
}
