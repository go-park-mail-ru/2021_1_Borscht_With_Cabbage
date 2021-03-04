package auth

import (
	"backend/api"
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
// true - такая сессия существует, false - нет
func checkSession(sessionToCheck string) (bool, error) {
	f, err := os.Open("sessions.txt")
	if err != nil {
		fmt.Println("error: ", err)
		return false, err
	}
	defer f.Close()
	buf := bufio.NewScanner(f)
	for {
		if !buf.Scan() {
			break
		}
		sessionJson := []byte(buf.Text())

		session := &api.Session{}
		err := json.Unmarshal(sessionJson, session)
		if err != nil {
			fmt.Println("error: ", err)
			return false, err
		}
		if session.Session == sessionToCheck {
			return true, nil
		}
	}
	return false, nil
}

// создание сессии для пользователя и привязка ее к пользователю(сейчас - по номеру телефону, в бд будет primary key)
// возвращает саму сессию чтобы вернуть ее на фронт
func createSession(number string) (string, error) {
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

	f, err := os.OpenFile("sessions.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}

	// создаем структуру чтобы замаршаллить и записать с файл
	sessionStruct := api.Session{
		Number:  number,
		Session: session,
	}
	sessionToWrite, err := json.Marshal(sessionStruct)
	if err != nil {
		return "", err
	}
	_, err = f.Write(sessionToWrite)
	if err != nil {
		fmt.Println("error: ", err)
		return "", err
	}
	return session, nil
}
