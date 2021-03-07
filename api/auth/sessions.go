package auth

import (
	"math/rand"
)

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
// true - такая сессия существует, false - нет
func checkSession(sessionToCheck string) (bool, error) {
	for _, session := range Sessions {
		if session.Session == sessionToCheck {
			return true, nil
		}
	}

	return false, nil
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

		isItExists, err := checkSession(session) // далее в цикле - проверка на уникальность
		if err != nil {
			return "", err
		}
		if !isItExists {
			break
		}
	}

	return session, nil
}
