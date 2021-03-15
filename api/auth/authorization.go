package auth

import (
	"backend/api"
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
	cook, err := cc.Cookie("borscht_session")
	if err != nil {
		return c.String(http.StatusOK, "error with request data")
	}

	fmt.Println(*cc.Sessions)
	session := cook.Value

	_, ok := (*cc.Sessions)[session]
	if ok {
		delete(*cc.Sessions, session)
	}

	fmt.Println(*cc.Sessions)

	DeleteResponseCookie(c)

	return c.JSON(http.StatusOK, "")
}

// handler авторизации
func LoginUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserAuth)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusUnauthorized, "error with request data")
	}

	for _, user := range *cc.Users {
		if (user.Email == newUser.Login || user.Phone == newUser.Login) && user.Password == newUser.Password {
			session := CreateSession(cc)

			(*cc.Sessions)[session] = user.Phone

			// далее - чтобы после авторизации пользователь перешел на главную
			SetResponseCookie(c, session)

			response := successResponse{user.Name, user.Avatar}
			return c.JSON(http.StatusOK, response)
		}
	}

	return c.JSON(http.StatusOK, "")
}

// handler регистрации
func CreateUser(c echo.Context) error {
	cc := c.(*api.CustomContext)
	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusUnauthorized, "error with request data")
	}

	for _, user := range *cc.Users {
		if (user.Phone == newUser.Phone) && user.Password == newUser.Password {
			return c.String(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
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
	return c.JSON(http.StatusOK, response)
}

func CheckAuth(c echo.Context) error {
	cc := c.(*api.CustomContext)

	user, err := GetUser(cc)
	if err != nil {
		return c.String(http.StatusUnauthorized, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}
