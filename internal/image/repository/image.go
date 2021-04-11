package repository

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/borscht/backend/internal/image"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
)

type imageRepo struct {
}

func NewImageRepo() image.ImageRepo {
	return &imageRepo{}
}

func (i imageRepo) UploadImage(ctx context.Context, filename string, image *multipart.FileHeader) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}

	return nil
}

func (i imageRepo) DeleteImage(ctx context.Context, filename string) error {
	err := os.Remove(filename)

	if err != nil {
		custError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, custError)
		return custError
	}

	logger.RepoLevel().InfoLog(ctx, logger.Fields{"remove file": filename})
	return nil
}
