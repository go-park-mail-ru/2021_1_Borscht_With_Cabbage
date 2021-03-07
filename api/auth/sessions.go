package auth

import (
	"backend/api"
	"math/rand"
)

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string, context *api.CustomContext) (string, bool) {
	session, isItExists := (*context.Sessions)["route"]
	if !isItExists {
		return "", false
	}
	return session, true
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
