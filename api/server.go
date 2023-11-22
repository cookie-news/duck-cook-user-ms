package api

import (
	"duck-cook-user-ms/controller"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Server struct {
	controller controller.Controller
}

func NewServer(controller controller.Controller) *Server {
	return &Server{controller}
}

func (server *Server) Start(port string) error {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWEDS_DOMAIN"))
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	})

	r.Use(func(ctx *gin.Context) {
		switch ctx.FullPath() {
		case
			"/v1/customer/:fieldName/:value",
			"/v1/customer":
			ctx.Next()
			return
		}

		auth := ctx.GetHeader("authorization")

		client := resty.New()
		client.BaseURL = os.Getenv("URL_AUTH")

		resp, _ := client.R().
			SetHeader("authorization", auth).
			Post("/v1/auth/verify-jwt")

		if resp.StatusCode() == http.StatusNoContent {
			ctx.Next()
			return
		} else {
			ctx.String(resp.StatusCode(), resp.String())
			ctx.Abort()
			return
		}
	})

	v1 := r.Group("/v1")
	{
		customer := v1.Group("/customer")
		{
			customer.GET("", server.controller.ListCustomersHandle)
			customer.GET("/:fieldName/:value", server.controller.GetUserByFieldHandle)
			customer.POST("", server.controller.CreateCustomerHandle)
		}
	}

	return r.Run(":" + port)
}
