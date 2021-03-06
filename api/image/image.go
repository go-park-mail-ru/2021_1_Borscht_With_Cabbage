package image

import (
	"backend/api"
	"fmt"
	"github.com/labstack/echo/v4"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// TODO обработка ошибок

func DownloadAvatar(c echo.Context) error {
	// TODO по сессии находим пользователя

	user, err := getUser(c)
	if err != nil {
		return err
	}

	filename := user.Avatar

	return c.File("static/avatar/" + filename)
}

func UploadAvatar(c echo.Context) error {
	// TODO проверка сессии

	// Читаем файл из пришедшего запроса
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// парсим расширение
	expansion := strings.Split(file.Header["Content-Type"][0], "/")[1]

	uid := getUniqId(file.Filename)
	filename := uid + "." + expansion
	err = setAvatarUser(c, filename)
	if err != nil {
		return err
	}

	// создаем файл у себя
	dst, err := os.Create("static/avatar/" + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// TODO return Json
	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename))

}

var user api.User

func getUser(c echo.Context) (api.User, error) {
	// TODO по сессии находим пользователя
	return user, nil
}

func setAvatarUser(c echo.Context, avatar string) error {
	user.Avatar = avatar

	return nil
}

func getUniqId(filename string) string {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	h := fnv.New32a()
	_, _ = h.Write([]byte(filename + hashingSalt))
	return strconv.Itoa(int(h.Sum32()))
}
