package auth_repository

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

type AuthRepository interface {
	CreateUser(user, pass, email string) (err error)
}

type authRepositoryImpl struct {
	clientHttp *resty.Client
}

func (repository authRepositoryImpl) CreateUser(user, pass, email string) (err error) {
	resp, err := repository.clientHttp.R().SetBody(map[string]any{
		"user":  user,
		"email": email,
		"pass":  pass,
	}).
		Post("/v1/customer")

	if resp.StatusCode() != http.StatusCreated {
		err = errors.New("não foi possivél criar o usuário")
	}

	return
}

func NewAuthRepository() AuthRepository {
	clientHttp := resty.New()
	clientHttp.BaseURL = os.Getenv("URL_AUTH")

	return &authRepositoryImpl{
		clientHttp,
	}
}
