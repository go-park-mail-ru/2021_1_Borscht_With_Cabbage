package auth

import (
	"backend/api"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserReg struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var Users []api.User

// данные о пользователях записаны в файлике, тут он парсится
func ParseUsers() error {
	users := make([]api.User, 0, 0)
	f, err := os.Open("users.txt")
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	defer f.Close()
	buf := bufio.NewScanner(f)
	for {
		if !buf.Scan() {
			break
		}
		userJson := []byte(buf.Text())

		user := &api.User{}
		err := json.Unmarshal(userJson, user)
		if err != nil {
			fmt.Println("error: ", err)
			return err
		}
		users = append(users, *user)
	}
	Users = users
	return nil
}

// записываем нового юзера в файлик при регистрации
func writeNewUser(userToRegister api.User) error {
	var userToWrite interface{} = userToRegister
	result, err := json.Marshal(userToWrite)
	if err != nil {
		return err
	}
	f, err := os.OpenFile("users.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	_, err = f.Write(result)
	if err != nil {
		fmt.Println("error: ", err)
		return err
	}
	return nil
}

// handler авторизации
func LogUser(c echo.Context) error {
	if err := ParseUsers(); err != nil {
		return err
	}
	newUser := new(UserAuth)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	for _, user := range Users {
		if (user.Name == newUser.Login || user.Number == newUser.Login) && user.Password == newUser.Password {
			// TODO тут должно быть обращение к функции, которая отдает json для главной страницы,
			// чтобы после авторизации пользователь перешел на главную
			return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
		}
	}
	return c.String(http.StatusUnauthorized, "")
}

// handler регистрации
func CreateUser(c echo.Context) error {
	if err := ParseUsers(); err != nil {
		return err
	}
	newUser := new(UserReg)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	for _, user := range Users {
		if (user.Number == newUser.Number) && user.Password == newUser.Password {
			return c.String(http.StatusUnauthorized, "user with this number already exists") // такой юзер уже есть
		}
	}

	userToRegister := api.User{
		Name:     newUser.Name,
		Email:    "",
		Password: newUser.Password,
		Number:   newUser.Number,
	}
	Users = append(Users, userToRegister)

	// записываем нового в файлик
	err := writeNewUser(userToRegister)
	if err != nil {
		return err
	}

	// TODO тут должно быть обращение к функции, которая отдает json для главной страницы,
	return c.String(http.StatusOK, "вместо этого текста тут json для формирования главной")
}

