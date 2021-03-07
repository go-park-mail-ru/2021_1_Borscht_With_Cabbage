package image

import (
	"backend/api"
	"backend/api/auth"
	"github.com/labstack/echo/v4"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func DownloadAvatar(c echo.Context) error {
	user, err := auth.GetUser(c.(*api.CustomContext))
	if err != nil {
		return getErrorJson(c, err)
	}

	filename := user.Avatar

	return getErrorJson(c, c.File("static/avatar/"+filename))
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
	expansion := filepath.Ext(file.Filename)

	uid, err := getUniqId(file.Filename)
	if err != nil {
		return getErrorJson(c, err)
	}

	filename := uid + expansion
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

	return DownloadAvatar(c)
}

func setAvatarUser(c echo.Context, avatar string) error {
	cc := c.(*api.CustomContext)
	currentUser, err := auth.GetUser(cc)
	if err != nil {
		return err
	}

	for i, user := range *cc.Users {
		if user == currentUser {
			user.Avatar = avatar
			(*cc.Users)[i] = user
		}
	}

	return nil
}

func getUniqId(filename string) (string, error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	// создаем хеш от названия файла
	hash := fnv.New32a()
	_, err := hash.Write([]byte(filename + hashingSalt))
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(hash.Sum32())), nil
}

func getErrorJson(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	return c.JSON(http.StatusOK, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
