package api

import (
	"duck-cook-user-ms/controllers"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Start(port string, client *mongo.Client) error {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWEDS_DOMAIN"))
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	})

	v1 := r.Group("/v1")
	{
		users := v1.Group("/customer")
		{
			users.GET("/", controllers.ListUsers)
			users.POST("/", func(ctx *gin.Context) { controllers.CreatCustomer(ctx, client) })
		}
	}

	return r.Run(":" + port)
}
