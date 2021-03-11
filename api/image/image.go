package image

import (
	"backend/api"
	"backend/api/auth"
	"backend/models"
	"github.com/labstack/echo/v4"
	"hash/fnv"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

const (
	HeadAvatar = "static/avatar/"
)

func UploadAvatar(c echo.Context) (string, models.ServerError) {
	// Читаем файл из пришедшего запроса
	file, err := c.FormFile("avatar")
	if err != nil {
		return "", models.BadRequest(err)
	}
	src, err := file.Open()
	if err != nil {
		return "", models.FailServer(err)
	}
	defer src.Close()

	// парсим расширение
	expansion := filepath.Ext(file.Filename)

	uid, localErr := getUniqId(file.Filename)
	if localErr != nil {
		return "", localErr
	}

	filename := HeadAvatar + uid + expansion

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return "", models.FailServer(err)
	}
	defer dst.Close()

	filename = api.Repository + filename
	localErr = setAvatarUser(c, filename)
	if localErr != nil {
		return "", localErr
	}

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return "", models.FailServer(err)
	}

	return filename, nil
}

func setAvatarUser(c echo.Context, avatar string) models.Error {
	cc := c.(*api.CustomContext)
	currentUser, err := auth.GetUser(cc)
	if err != nil {
		return models.AuthorizationErr(err)
	}

	for i, user := range *cc.Users {
		if user == currentUser {
			user.Avatar = avatar
			(*cc.Users)[i] = user
		}
	}

	return nil
}

func getUniqId(filename string) (string, models.Error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	// создаем хеш от названия файла
	hash := fnv.New32a()
	_, err := hash.Write([]byte(filename + hashingSalt))
	if err != nil {
		return "", models.FailServer(err)
	}

	return strconv.Itoa(int(hash.Sum32())), nil
}

