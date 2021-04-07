package image

import (
	"context"
	"mime/multipart"
)

type ImageRepo interface {
	UploadImage(ctx context.Context, filename string, image *multipart.FileHeader) error
	DeleteImage(ctx context.Context, filename string) error
}
