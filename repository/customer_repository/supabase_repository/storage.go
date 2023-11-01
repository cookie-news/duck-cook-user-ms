package supabase_repository

import (
	"duck-cook-user-ms/api/repository"
	"mime/multipart"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

type storageImpl struct {
	client storage_go.Client
}

func (c *storageImpl) UploadImage(image multipart.FileHeader, user string) (string, error) {

	file, err := image.Open()
	contentType := image.Header.Get("Content-Type")

	if err != nil {
		return "", err
	}

	result, err := c.client.UploadFile(os.Getenv("SUPABASE_BUCKET_ID"), user, file, storage_go.FileOptions{ContentType: &contentType})

	if err != nil {
		return "", err
	}

	return result.Key, err
}

func (c *storageImpl) GetPublicUrl(filename string) string {

	result := c.client.GetPublicUrl(os.Getenv("SUPABASE_BUCKET_ID"), filename)
	return result.SignedURL
}

func New(client storage_go.Client) repository.CustomerStorage {
	return &storageImpl{client}
}
