package auth

import (
	"backend/api/domain"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func SetResponseCookie(c echo.Context, session string) {
	sessionCookie := http.Cookie {
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     domain.SessionCookie,
		Value:    session,
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

func DeleteResponseCookie(c echo.Context) {
	sessionCookie := http.Cookie {
		Expires:  time.Now().Add(-24 * time.Hour),
		Name:     domain.SessionCookie,
		Value:    "session",
		HttpOnly: true,
	}
	c.SetCookie(&sessionCookie)
}

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string, context *domain.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание уникальной сессии
func CreateSession(context *domain.CustomContext) string {
	session := ""

	for {
		session = uuid.New().String()

		_, isItExists := CheckSession(session, context) // далее в цикле - проверка на уникальность
		if isItExists == false {                        // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	return session
}

func GetUser(context *domain.CustomContext) (domain.User, error) {
	sessionError := errors.New("session error")
	session, err := context.Cookie(domain.SessionCookie)
	if err != nil {
		return domain.User{}, sessionError
	}

	phone, ok := CheckSession(session.Value, context)
	if !ok {
		return domain.User{}, sessionError
	}

	for _, user := range *context.Users {
		if user.Phone == phone {
			return user, nil
		}
	}

	return domain.User{}, errors.New("user not found")
}
