package http

import (
	"backend/api/domain"
	"backend/api/user/delivery/http/middleware"
	errors "backend/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserHandler struct {
	UserUcase    domain.UserUsecase
	SessionUcase domain.SessionUsecase
	ImageUcase domain.ImageUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewUserHandler(e *echo.Echo,
					uus domain.UserUsecase,
					sus domain.SessionUsecase,
					iuc domain.ImageUsecase) {
	handler := &UserHandler{
		UserUcase: uus,
		SessionUcase: sus,
		ImageUcase: iuc,
	}

	initMiddleware := middleware.InitMiddleware(uus, sus)
	e.Use(initMiddleware.Auth)

	e.POST("/signin", handler.LoginUser)
	e.POST("/signup", handler.Create)
	e.GET("/user", handler.GetUserData)
	e.PUT("/user", handler.EditProfile)
	e.GET("/auth", handler.CheckAuth)
	e.GET("/logout", handler.LogoutUser)
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

func deleteResponseCookie(c echo.Context) {
	sessionCookie := http.Cookie {
		Expires:  time.Now().Add(-24 * time.Hour),
		Name:     domain.SessionCookie,
		Value:    "session",
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

func (h *UserHandler) Create(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	newUser := new(domain.UserReg)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	userToRegister := domain.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Phone:    newUser.Phone,
		Avatar:   domain.DefaultAvatar,
	}

	if err := h.UserUcase.Create(cc, userToRegister); err != nil {
		return cc.SendERR(err)
	}

	session, err := h.SessionUcase.Create(cc, newUser.Phone)
	if err != nil {
		return cc.SendERR(err)
	}

	setResponseCookie(c, session)

	response := domain.SuccessResponse{Name: userToRegister.Name}
	return cc.SendOK(response)
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	cc := c.(*domain.CustomContext)
	newUser := new(domain.UserAuth)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	oldUser, err := h.UserUcase.GetByLogin(cc, *newUser)
	if err != nil {
		return cc.SendERR(err)
	}

	session, err := h.SessionUcase.Create(cc, oldUser.Phone)
	if err != nil {
		return cc.SendERR(err)
	}

	setResponseCookie(c, session)

	response := domain.SuccessResponse{Name: oldUser.Name, Avatar: oldUser.Avatar}
	return cc.SendOK(response)
}

func (h *UserHandler) GetUserData(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	if cc.User == nil {
		userError := errors.Authorization("not authorization")
		return cc.SendERR(userError)
	}

	return cc.SendOK(*cc.User)
}

func (h *UserHandler) EditProfile(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	// TODO убрать часть логики в usecase
	profileEdits := new(domain.UserData)
	formParams, err := c.FormParams()
	if err != nil {
		return errors.Create(http.StatusBadRequest, "invalid data form")
	}

	profileEdits.Name = formParams.Get("name")
	profileEdits.Phone = formParams.Get("number")
	profileEdits.Email = formParams.Get("email")
	profileEdits.Password = formParams.Get("password")
	profileEdits.PasswordOld = formParams.Get("password_current")
	fmt.Println(profileEdits)

	file, err := c.FormFile("avatar")
	if err != nil {
		return cc.SendERR(errors.BadRequest(err.Error()))
	}
	filename, err := h.ImageUcase.UploadAvatar(cc, file)
	if err != nil {
		return cc.SendERR(err)
	}

	profileEdits.Avatar = filename

	if cc.User == nil {
		userError := errors.Authorization("not authorization")
		return cc.SendERR(userError)
	}

	err = h.UserUcase.Update(cc, *profileEdits)
	if err != nil {
		return cc.SendERR(err)
	}
	err = h.SessionUcase.UpdateValue(cc, profileEdits.Phone, cc.User.Phone)
	if err != nil {
		return cc.SendERR(err)
	}

	return cc.SendOK(profileEdits)
}

func (h *UserHandler) CheckAuth(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	if cc.User == nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	return cc.SendOK(cc.User)
}

func (h *UserHandler) LogoutUser(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	cook, err := cc.Cookie(domain.SessionCookie)
	if err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	h.SessionUcase.Delete(cc, cook.Value)

	deleteResponseCookie(c)

	return cc.SendOK("")
}
