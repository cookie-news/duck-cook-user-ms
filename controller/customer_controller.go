package controller

import (
	"duck-cook-user-ms/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateCustomerHandle(ctx *gin.Context) {
	var customer entity.Customer
	if err := ctx.ShouldBindJSON(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao decodificar o JSON"})
		return
	}

	_, err := c.customerUsecase.CreateCustomer(customer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"customer": customer,
	})
}

func (c *Controller) ListCustomersHandle(ctx *gin.Context) {
	customers, err := c.customerUsecase.ListCustomers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customers)
}
