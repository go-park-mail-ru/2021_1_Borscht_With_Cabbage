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
	errors "github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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

func (u *userUsecase) Create(ctx context.Context, newUser models.User) (*models.User, error) {

	// TODO валидация какая нибудь
	newUser.Avatar = config.DefaultAvatar

	uid, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.Uid = uid
	return &newUser, nil
}

func (u *userUsecase) CheckUserExists(ctx context.Context, user models.UserAuth) (*models.User, error) {
	return u.userRepository.CheckUserExists(ctx, user)
}

func (u *userUsecase) GetByUid(ctx context.Context, uid int) (models.User, error) {
	return u.userRepository.GetByUid(ctx, uid)
}

func (u *userUsecase) Update(ctx context.Context, newUser models.UserData, uid int) error {
	// TODO валидация

	return u.userRepository.Update(ctx, newUser, uid)
}

func (u *userUsecase) UploadAvatar(ctx context.Context, image *multipart.FileHeader) (string, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := getUniqId(ctx, image.Filename)
	if localErr != nil {
		return "", localErr
	}

	filename := HeadAvatar + uid + expansion
	if err := u.userRepository.UploadAvatar(ctx, image, filename); err != nil {
		return "", err
	}

	return config.Repository + filename, nil
}

func getUniqId(ctx context.Context, filename string) (string, error) {
	// создаем рандомную последовательность чтобы точно названия не повторялись
	hashingSalt := strconv.Itoa(rand.Int() % 1000)

	// создаем хеш от названия файла
	hash := fnv.New32a()
	_, err := hash.Write([]byte(filename + hashingSalt))
	if err != nil {
		custErr := errors.FailServerError(err.Error())
		logger.UsecaseLevel().ErrorLog(ctx, custErr)
		return "", custErr
	}

	return strconv.Itoa(int(hash.Sum32())), nil
}
