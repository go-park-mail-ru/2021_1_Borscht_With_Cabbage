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

const (
	HeadAvatar = "static/avatar/"
	Repository = "http://89.208.197.150:5000/"
)

func DownloadAvatar(c echo.Context) error {
	//user, err := auth.GetUser(c.(*api.CustomContext))
	//fmt.Println("hi")
	//if err != nil {
	//	fmt.Println("err")
	//	return getErrorJson(c, err)
	//}
	//
	//fmt.Println(user)
	//
	//filename := user.Avatar
	//
	//return getErrorJson(c, c.File("static/avatar/"+filename))
	return getErrorJson(c, c.File("static/avatar/stas.jpg"))
}

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

	uid, err := getUniqId(file.Filename)
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

	filename = Repository + filename
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
