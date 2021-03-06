package image

import (
	"backend/api"
	"github.com/labstack/echo/v4"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DownloadAvatar(c echo.Context) error {
	user, err := getUser(c)
	if err != nil {
		return getErrorJson(c, err)
	}

	filename := user.Avatar

	return getErrorJson(c, c.File("static/avatar/" + filename))
}

func UploadAvatar(c echo.Context) error {
	// Читаем файл из пришедшего запроса
	file, err := c.FormFile("avatar")
	if err != nil {
		return getErrorJson(c, err)
	}
	src, err := file.Open()
	if err != nil {
		return getErrorJson(c, err)
	}
	defer src.Close()

	// парсим расширение
	expansion := strings.Split(file.Header["Content-Type"][0], "/")[1]

	uid, err := getUniqId(file.Filename)
	if err != nil {
		return getErrorJson(c, err)
	}

	filename := uid + "." + expansion
	err = setAvatarUser(c, filename)
	if err != nil {
		return getErrorJson(c, err)
	}

	// создаем файл у себя
	dst, err := os.Create("static/avatar/" + filename)
	if err != nil {
		return getErrorJson(c, err)
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return getErrorJson(c, err)
	}

	return c.JSON(http.StatusOK, struct{
		Filename string `json:"filename"`
	}{
		Filename: filename,
	})
}

var user api.User

func getUser(c echo.Context) (api.User, error) {
	// TODO по сессии находим пользователя
	return user, nil
}

func setAvatarUser(c echo.Context, avatar string) error {
	// TODO проверка сессии
	user.Avatar = avatar

	return nil
}

func getUniqId(filename string) (string, error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	h := fnv.New32a()
	_, err := h.Write([]byte(filename + hashingSalt))
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(h.Sum32())), nil
}

func getErrorJson(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	return c.JSON(http.StatusOK, struct{
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
