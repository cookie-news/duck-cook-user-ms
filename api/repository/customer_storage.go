package repository

import (
	"mime/multipart"
)

type CustomerStorage interface {
	UploadImage(image multipart.FileHeader, user string) (string, error)
	GetPublicUrl(filename string) string
}
