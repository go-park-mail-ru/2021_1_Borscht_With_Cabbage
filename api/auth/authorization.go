package auth

import (
	"backend/api"
	errors "backend/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

func LogoutUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	cook, err := cc.Cookie(api.SessionCookie)
	if err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	fmt.Println(*cc.Sessions)
	session := cook.Value

	_, ok := (*cc.Sessions)[session]
	if ok {
		delete(*cc.Sessions, session)
	}

	fmt.Println(*cc.Sessions)

	DeleteResponseCookie(c)

	return cc.SendOK("")
}

// handler авторизации
func LoginUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserAuth)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	for _, user := range *cc.Users {
		if (user.Email == newUser.Login || user.Phone == newUser.Login) && user.Password == newUser.Password {
			session := CreateSession(cc)

			(*cc.Sessions)[session] = user.Phone

			// далее - чтобы после авторизации пользователь перешел на главную
			SetResponseCookie(c, session)

			response := successResponse{user.Name, user.Avatar}
			return cc.SendOK(response)
		}
	}

	return cc.SendOK("")
}

// handler регистрации
func CreateUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	for _, user := range *cc.Users {
		if (user.Phone == newUser.Phone) && user.Password == newUser.Password {
			sendErr := errors.Create(http.StatusUnauthorized, "user with this number already exists")
			return cc.SendERR(sendErr) // такой юзер уже есть
		}
	}

	userToRegister := api.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Phone:    newUser.Phone,
		Avatar:   api.DefaultAvatar,
	}

	// записываем нового
	*cc.Users = append(*cc.Users, userToRegister)

	session := CreateSession(cc)
	(*cc.Sessions)[session] = newUser.Phone

	// далее - чтобы после авторизации пользователь перешел на главную
	SetResponseCookie(c, session)

	response := successResponse{userToRegister.Name, ""}
	return cc.SendOK(response)
}

func CheckAuth(c echo.Context) error {
	cc := c.(*api.CustomContext)

	user, err := GetUser(cc)
	if err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	return cc.SendOK(user)
}
