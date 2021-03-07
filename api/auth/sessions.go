package auth

import (
	"backend/api"
	"math/rand"
)

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string) string {
	for _, session := range api.Sessions {
		if session.Session == sessionToCheck {
			return session.Number
		}
	}

	return ""
}

// создание сессии для пользователя и привязка ее к пользователю(сейчас - по номеру телефону, в бд будет primary key)
// возвращает саму сессию чтобы вернуть ее на фронт
func createSession() (string, error) {
	session := ""

	for {
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		sessionMaking := make([]rune, sessionLen)

		for i := range sessionMaking {
			sessionMaking[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		session = string(sessionMaking)

		isItExists := CheckSession(session) // далее в цикле - проверка на уникальность
		if isItExists == "" { // не получили привязанного к сессии пользователя, следовательно, не существует
			break
		}
	}

	return session, nil
}
