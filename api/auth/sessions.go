package auth

import (
	"backend/api"
	"errors"
	"math/rand"
)

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string, context *api.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)["route"]
	if !isItExists {
		return "", false
	}
	return number, true
}

// создание сессии для пользователя и привязка ее к пользователю(сейчас - по номеру телефону, в бд будет primary key)
// возвращает саму сессию чтобы вернуть ее на фронт
func createSession(context *api.CustomContext) string {
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
	session, err := context.Cookie("session")
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
