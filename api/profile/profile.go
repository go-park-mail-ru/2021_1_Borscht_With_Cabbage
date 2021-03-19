package profile

import (
	"backend/api/auth"
	"backend/api/domain"
	"backend/api/image"
	"backend/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserData struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"number"`
	Password       string `json:"password"`
	PasswordOld    string `json:"password_current"`
	Avatar         string `json:"avatar"`
}

// отдать данные о юзере, чтобы загрузить профиль
func GetUserData(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	user, err := auth.GetUser(cc)
	if err != nil {
		return cc.SendERR(err)
	}

	return cc.SendOK(user)
}

// сохранить изменения в профиле
func EditProfile(c echo.Context) error {
	cc := c.(*domain.CustomContext)

	profileEdits := new(UserData)
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

	srcFile, err := image.UploadAvatar(c)

	profileEdits.Avatar = srcFile

	// по куке находим, что за юзер
	session, err := cc.Cookie("borscht_session")
	if err != nil {
		return cc.SendERR(errors.Authorization("user not found"))
	}
	number, isItExists := auth.CheckSession(session.Value, cc)
	if isItExists == false {
		return cc.SendERR(errors.Authorization("user not found"))
	}

	for i, user := range *cc.Users {
		if user.Email == profileEdits.Email && user.Phone != number { // если у кого-то другого уже есть такой email
			return cc.SendERR(errors.Create(http.StatusBadRequest, "user with this email already exists"))
		}
		if user.Phone == profileEdits.Phone && user.Phone != number { // если у кого-то другого уже есть такой телефон
			return cc.SendERR(errors.Create(http.StatusBadRequest, "user with this number already exists"))
		}

		if user.Phone == number {
			if profileEdits.Password != "" {
				if profileEdits.PasswordOld != user.Password {
					fmt.Println(profileEdits.PasswordOld, " ", user.Password)
					return cc.SendERR(errors.Create(http.StatusBadRequest, "invalid old password"))
				}
				(*cc.Users)[i].Password = profileEdits.Password
			}

			(*cc.Users)[i].Phone = profileEdits.Phone
			(*cc.Users)[i].Email = profileEdits.Email
			(*cc.Users)[i].Name = profileEdits.Name

			for j, numSession := range *cc.Sessions {
				if numSession == number {
					(*cc.Sessions)[j] = profileEdits.Phone
				}
			}

			fmt.Println(*cc.Users)
			return cc.SendOK(profileEdits)
		}
	}

	return cc.SendERR(errors.Authorization("user not found"))
}
