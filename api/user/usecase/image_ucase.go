package usecase

import (
	"backend/api/domain"
	"backend/api/domain/user"
	errors "backend/utils"
	"hash/fnv"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"strconv"
)

const (
	HeadAvatar = "static/avatar/"
)

type imageUsecase struct {
	imageRepository user.ImageRepo
}

func NewImageUsecase(repo user.ImageRepo) user.ImageUsecase {
	return &imageUsecase{
		imageRepository: repo,
	}
}

func (i imageUsecase) UploadAvatar(ctx *domain.CustomContext, image *multipart.FileHeader) (string, error) {
	// парсим расширение
	expansion := filepath.Ext(image.Filename)

	uid, localErr := getUniqId(image.Filename)
	if localErr != nil {
		return "", localErr
	}

	filename := HeadAvatar + uid + expansion
	if err := i.imageRepository.UploadAvatar(image, filename); err != nil {
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
