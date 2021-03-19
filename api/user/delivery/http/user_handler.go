package http

import (
	"backend/api/auth"
	"backend/api/domain"
	"backend/api/profile"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserHandler struct {
	UserUcase domain.UserUsecase
	SessionUcase domain.SessionUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewUserHandler(e *echo.Echo, uus domain.UserUsecase, sus domain.SessionUsecase) {
	handler := &UserHandler{
		UserUcase: uus,
		SessionUcase: sus,
	}

	e.POST("/signin", auth.LoginUser)
	e.POST("/signup", handler.Create)
	e.GET("/user", profile.GetUserData)
	e.PUT("/user", profile.EditProfile)
	e.GET("/auth", auth.CheckAuth)
	e.GET("/logout", auth.LogoutUser)
}

type UserReg struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type successResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func setResponseCookie(c echo.Context, session string) {
	sessionCookie := http.Cookie {
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     domain.SessionCookie,
		Value:    session,
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

func (h *UserHandler) Create(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusUnauthorized, "error with request data")
	}

	userToRegister := domain.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Phone:    newUser.Phone,
		Avatar:   domain.DefaultAvatar,
	}

	if err := h.UserUcase.Create(cc, userToRegister); err != nil {
		return err
	}

	session, err := h.SessionUcase.Create(cc, newUser.Phone)
	if err != nil {
		return err
	}

	setResponseCookie(c, session)

	response := successResponse{userToRegister.Name, ""}
	return c.JSON(http.StatusOK, response)
}
