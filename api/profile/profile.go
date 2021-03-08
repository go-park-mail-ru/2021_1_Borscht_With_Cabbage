package profile

import (
	"backend/api"
	"backend/api/auth"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"number"`
}

// отдать данные о юзере, чтобы загрузить профиль
func GetUserData(c echo.Context) error {
	cc := c.(*api.CustomContext)

	user, err := auth.GetUser(cc)
	if err != nil {
		return err
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(userData))
}

// сохранить изменения в профиле
func EditProfile(c echo.Context) error {
	cc := c.(*api.CustomContext)

	profileEdits := new(UserData)
	if err := cc.Bind(profileEdits); err != nil {
		return cc.String(http.StatusUnauthorized, "error with request data")
	}

	// по куке находим, что за юзер
	session, err := cc.Cookie("session")
	if err != nil {
		return c.String(http.StatusUnauthorized, "no session found")
	}
	number, isItExists := auth.CheckSession(session.Value, cc)
	if isItExists == false {
		return c.String(http.StatusUnauthorized, "user not found")
	}

	for _, user := range *cc.Users {
		if user.Email == profileEdits.Email && user.Phone != number { // если у кого-то другого уже есть такой email
			return c.String(http.StatusNotImplemented, "user with this email already exists")
		}
		if user.Phone == profileEdits.Phone && user.Phone != number { // если у кого-то другого уже есть такой телефон
			return c.String(http.StatusNotImplemented, "user with this phone number already exists")
		}

		if user.Phone == number {
			user.Phone = profileEdits.Phone
			user.Email = profileEdits.Email
			user.Name = profileEdits.Name
		}
	}

	return c.String(http.StatusUnauthorized, "user not found")
}
