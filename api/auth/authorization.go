package auth

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"net/http"
)

var sessionLen = 15

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserReg struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// handler авторизации
func LoginUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserAuth)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusUnauthorized, "error with request data")
	}

	for _, user := range *cc.Users {
		if (user.Name == newUser.Login || user.Phone == newUser.Login) && user.Password == newUser.Password {
			session := createSession(cc)

			(*cc.Sessions)[session] = user.Phone
			// TODO тут должно быть обращение к функции, которая отдает json для главной страницы, и созданную выше сессию в том числе
			// чтобы после авторизации пользователь перешел на главную
			return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
		}
	}
	return c.String(http.StatusUnauthorized, "")
}

// handler регистрации
func CreateUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusUnauthorized, "error with request data")
	}

	for _, user := range *cc.Users {
		if (user.Phone == newUser.Number) && user.Password == newUser.Password {
			return c.String(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	userToRegister := api.User{
		Name:     newUser.Name,
		Email:    "",
		Password: newUser.Password,
		Phone:    newUser.Number,
	}

	// записываем нового
	*cc.Users = append(*cc.Users, userToRegister)

	session := createSession(cc)
	(*cc.Sessions)[session] = newUser.Number

	// TODO тут должно быть обращение к функции, которая отдает json для главной страницы,
	// и созданную выше сессию в том числе
	return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
}
