package main

import (
	"duck-cook-user-ms/api"
	"duck-cook-user-ms/controller"
	"duck-cook-user-ms/pkg/mongodb"
	"duck-cook-user-ms/pkg/supabase"
	"duck-cook-user-ms/repository/auth_repository"
	mongodb_repository "duck-cook-user-ms/repository/customer_repository/mongo_repository"
	"duck-cook-user-ms/repository/customer_repository/supabase_repository"

	"duck-cook-user-ms/usecase"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Server api.Server
}

func NewAppConfig() AppConfig {
	_ = godotenv.Load()

	mongoDb := mongodb.ConnectMongo()
	supabase := supabase.ConnectStorage()

	repositoryCustomer := mongodb_repository.New(mongoDb)
	storageCustomer := supabase_repository.New(supabase)
	authRepository := auth_repository.NewAuthRepository()

	customerUsecase := usecase.NewCustomerUseCase(repositoryCustomer, storageCustomer, authRepository)

	controller := controller.NewController(customerUsecase)
	server := api.NewServer(controller)

	return AppConfig{
		Server: *server,
	}
}
