package usecase

import (
	"context"
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

func (u *userUsecase) Create(ctx context.Context, newUser models.User) error {
	// TODO валидация какая нибудь
	newUser.Avatar = config.DefaultAvatar

	return u.userRepository.Create(ctx, newUser)
}

func (u *userUsecase) GetByLogin(ctx context.Context, user models.UserAuth) (models.User, error) {
	return u.userRepository.GetByLogin(ctx, user)
}

func (u *userUsecase) GetByNumber(ctx context.Context, number string) (models.User, error) {
	return u.userRepository.GetByNumber(ctx, number)
}

func (u *userUsecase) Update(ctx context.Context, newUser models.UserData) error {
	// TODO валидация

	return u.userRepository.Update(ctx, newUser)
}

func (u *userUsecase) UploadAvatar(ctx context.Context, image *multipart.FileHeader) (string, error) {
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

	return filename, nil
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