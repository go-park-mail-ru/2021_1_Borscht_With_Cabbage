package user

import (
	"backend/api/domain"
	"mime/multipart"
)

type Image struct {
}

type ImageUsecase interface {
	UploadAvatar(ctx *domain.CustomContext, image *multipart.FileHeader) (string, error)
}

type ImageRepo interface {
	UploadAvatar(image *multipart.FileHeader, filename string) error
}
