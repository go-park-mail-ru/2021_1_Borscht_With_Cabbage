package domain

import (
	"mime/multipart"
)

type Image struct {
}

type ImageUsecase interface {
	UploadAvatar(ctx *CustomContext, image *multipart.FileHeader) (string, error)
}

type ImageRepo interface {
	UploadAvatar(image *multipart.FileHeader, filename string) error
}
