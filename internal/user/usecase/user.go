package usecase

import (
	"hash/fnv"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/borscht/backend/config"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	errors "github.com/borscht/backend/utils"
)

// TODO: хранить статику в /var/...

const (
	HeadAvatar = "static/avatar/"
)

type userUsecase struct {
	userRepository user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) user.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u *userUsecase) Create(newUser models.User) (*models.User, error) {

	// TODO валидация какая нибудь
	newUser.Avatar = config.DefaultAvatar

	uid, err := u.userRepository.Create(newUser)
	if err != nil {
		return nil, err
	}
	newUser.Uid = uid
	return &newUser, nil
}

func (u *userUsecase) CheckUserExists(user models.UserAuth) (*models.User, error) {
	return u.userRepository.CheckUserExists(user)
}

func (u *userUsecase) GetByUid(uid int) (models.User, error) {
	return u.userRepository.GetByUid(uid)
}

func (u *userUsecase) Update(newUser models.UserData, uid int) error {
	// TODO валидация

	return u.userRepository.Update(newUser, uid)
}

func (u *userUsecase) UploadAvatar(image *multipart.FileHeader) (string, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := getUniqId(image.Filename)
	if localErr != nil {
		return "", localErr
	}

	filename := HeadAvatar + uid + expansion
	if err := u.userRepository.UploadAvatar(image, filename); err != nil {
		return "", err
	}

	return config.Repository + filename, nil
}

func getUniqId(filename string) (string, error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	// создаем хеш от названия файла
	hash := fnv.New32a()
	_, err := hash.Write([]byte(filename + hashingSalt))
	if err != nil {
		return "", errors.FailServer(err.Error())
	}

	return strconv.Itoa(int(hash.Sum32())), nil
}
