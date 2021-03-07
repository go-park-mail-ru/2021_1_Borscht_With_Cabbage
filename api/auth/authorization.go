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

var Users []api.User
var Sessions []api.Session

// handler авторизации
func LogUser(c echo.Context) error {
	newUser := new(UserAuth)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	for _, user := range Users {
		if (user.Name == newUser.Login || user.Number == newUser.Login) && user.Password == newUser.Password {
			session, err := createSession()
			if err != nil {
				return err
			}

			sessionToRead := api.Session{Session: session, Number: user.Number}
			Sessions = append(Sessions, sessionToRead)
			// TODO тут должно быть обращение к функции, которая отдает json для главной страницы, и созданную выше сессию в том числе
			// чтобы после авторизации пользователь перешел на главную
			return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
		}
	}
	return c.String(http.StatusUnauthorized, "")
}

// handler регистрации
func CreateUser(c echo.Context) error {
	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	for _, user := range Users {
		if (user.Number == newUser.Number) && user.Password == newUser.Password {
			return c.String(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	userToRegister := api.User{
		Name:     newUser.Name,
		Email:    "",
		Password: newUser.Password,
		Number:   newUser.Number,
	}

	// записываем нового
	Users = append(Users, userToRegister)

	session, err := createSession()
	if err != nil {
		return err
	}
	sessionToRead := api.Session{Session: session, Number: newUser.Number}
	Sessions = append(Sessions, sessionToRead)

	// TODO тут должно быть обращение к функции, которая отдает json для главной страницы,
	// и созданную выше сессию в том числе
	return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
}
