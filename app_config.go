package main

import (
	"duck-cook-user-ms/api"
	"duck-cook-user-ms/controller"
	"duck-cook-user-ms/db"
	customerrepository "duck-cook-user-ms/repository/customer_repository"
	"duck-cook-user-ms/usecase"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server api.Server
}

func NewAppConfig() AppConfig {
	_ = godotenv.Load()

	mongoDb := db.ConnectMongo()

	repositoryCustomer := customerrepository.New(mongoDb)
	customerUsecase := usecase.NewCustomerUseCase(repositoryCustomer)

	controller := controller.NewController(customerUsecase)
	server := api.NewServer(controller)

	return AppConfig{
		Server: *server,
	}
}
