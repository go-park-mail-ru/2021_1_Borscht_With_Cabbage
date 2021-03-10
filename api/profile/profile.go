package profile

import (
	"backend/api"
	"backend/api/auth"
	"backend/api/image"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"number"`
	Avatar string `json:"avatar"`
}

// отдать данные о юзере, чтобы загрузить профиль
func GetUserData(c echo.Context) error {
	cc := c.(*api.CustomContext)

	user, err := auth.GetUser(cc)
	if err != nil {
		return c.String(http.StatusUnauthorized, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}

// сохранить изменения в профиле
func EditProfile(c echo.Context) error {
	cc := c.(*api.CustomContext)

	profileEdits := new(UserData)
	//if err := cc.Bind(profileEdits); err != nil {
	//	return cc.String(http.StatusUnauthorized, "error with request data")
	//}
	formParams, err := c.FormParams()
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusOK, err)
	}

	profileEdits.Name = formParams.Get("name")
	profileEdits.Phone = formParams.Get("number")
	profileEdits.Email = formParams.Get("email")

	fmt.Println(profileEdits)

	srcFile, err := image.UploadAvatar(c)

	profileEdits.Avatar = srcFile

	// по куке находим, что за юзер
	session, err := cc.Cookie("borscht_session")
	if err != nil {
		return c.String(http.StatusUnauthorized, "no session found")
	}
	number, isItExists := auth.CheckSession(session.Value, cc)
	if isItExists == false {
		return c.String(http.StatusUnauthorized, "user not found")
	}

	for i, user := range *cc.Users {
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
		(*cc.Users)[i] = user

		return c.JSON(http.StatusOK, profileEdits)
	}

	return c.String(http.StatusUnauthorized, "user not found")
}
