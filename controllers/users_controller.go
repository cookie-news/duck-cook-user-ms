package controllers

import (
	"context"
	"duck-cook-user-ms/db/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListUsers(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"id": "1", "name": "Lucas"})
}

func CreatCustomer(ctx *gin.Context, client *mongo.Client) {
	coll := client.Database("duckcook").Collection("users")

	newCustomer := model.Customer{
		Email:            "lucasjoao85@gmail.com",
		User:             "lucasnascimento",
		Pass:             "lucas123",
		Name:             "Lucas Nascimento",
		ImageProfilePath: "https://www.goole.com/imagem.png",
	}

	_, err := coll.InsertOne(context.TODO(), newCustomer)

	if err != nil {
		ctx.JSON(500, gin.H{"message": "Ocorreu um erro ao tentar criar o usuário"})
		panic(err)
	}

	ctx.JSON(200, gin.H{"message": "Usuario criado com sucesso"})

}
