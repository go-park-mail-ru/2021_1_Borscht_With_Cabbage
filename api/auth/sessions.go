package auth

import (
	"backend/api"
	"errors"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)

var sessionLen = 15

func SetResponseCookie(c echo.Context, session string) {
	sessionCookie := new(http.Cookie)
	sessionCookie.Expires = time.Now().Add(24 * time.Hour)
	sessionCookie.Name = "borscht_session"
	sessionCookie.Value = session
	sessionCookie.HttpOnly = true
	c.SetCookie(sessionCookie)
}

func DeleteResponseCookie(c echo.Context) {
	sessionCookie := new(http.Cookie)
	sessionCookie.Expires = time.Now().Add(-24 * time.Hour)
	sessionCookie.Name = "borscht_session"
	sessionCookie.Value = ""
	sessionCookie.HttpOnly = true
	c.SetCookie(sessionCookie)
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string, context *api.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание сессии для пользователя и привязка ее к пользователю(сейчас - по номеру телефону, в бд будет primary key)
// возвращает саму сессию чтобы вернуть ее на фронт
func CreateSession(context *api.CustomContext) string {
	session := ""

	for {
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		sessionMaking := make([]rune, sessionLen)

		for i := range sessionMaking {
			sessionMaking[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		session = string(sessionMaking)

		_, isItExists := CheckSession(session, context) // далее в цикле - проверка на уникальность
		if isItExists == false {                        // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	return session
}

func GetUser(context *api.CustomContext) (api.User, error) {
	sessionError := errors.New("session error")
	session, err := context.Cookie("borscht_session")
	if err != nil {
		return api.User{}, sessionError
	}

	phone, ok := CheckSession(session.Value, context)
	if !ok {
		return api.User{}, sessionError
	}

	for _, user := range *context.Users {
		if user.Phone == phone {
			return user, nil
		}
	}

	return api.User{}, errors.New("user not found")
}
