package controller

import "duck-cook-user-ms/usecase"

type Controller struct {
	customerUsecase usecase.CustomerUsecase
}

func NewController(
	customerUsecase usecase.CustomerUsecase,
) Controller {
	return Controller{
		customerUsecase,
	}
}
