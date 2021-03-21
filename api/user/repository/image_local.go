package repository

import (
	"backend/api/domain"
	errors "backend/utils"
	"io"
	"mime/multipart"
	"os"
)

type imageRepo struct {
}

func NewImageRepo() domain.ImageRepo {
	return &imageRepo{}
}

func (i imageRepo) UploadAvatar(image *multipart.FileHeader,
								filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		return errors.FailServer(err.Error())
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return errors.FailServer(err.Error())
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return errors.FailServer(err.Error())
	}

	return nil
}
