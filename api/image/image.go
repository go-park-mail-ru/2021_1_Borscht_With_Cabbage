package image

import (
	"backend/api/auth"
	"backend/api/domain"
	"github.com/labstack/echo/v4"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	HeadAvatar = "static/avatar/"
)

func UploadAvatar(c echo.Context) (string, error) {
	// Читаем файл из пришедшего запроса
	file, err := c.FormFile("avatar")
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// парсим расширение
	expansion := filepath.Ext(file.Filename)

	uid, err := getUniqFileId(file.Filename)
	if err != nil {
		return "", err
	}

	filename := HeadAvatar + uid + expansion

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	filename = domain.Repository + filename
	err = setAvatarUser(c, filename)
	if err != nil {
		return "", err
	}

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

func setAvatarUser(c echo.Context, avatar string) error {
	cc := c.(*domain.CustomContext)
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

func getUniqFileId(filename string) (string, error) {
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
