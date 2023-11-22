package controller

import (
	"duck-cook-user-ms/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CreateCustomerHandle(ctx *gin.Context) {
	var customer entity.Customer

	if err := ctx.ShouldBind(&customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerResponse, err := c.customerUsecase.CreateCustomer(customer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if customer.Image != nil {
		url, err := c.customerUsecase.UploadImage(*customer.Image, customer.User)

		if err != nil {
			c.customerUsecase.DeleteCustomer(customer.ID)
			ctx.JSON(http.StatusCreated, gin.H{
				"customer": customer,
				"error":    "An error occurred while saving your profile photo, but the username was created successfully",
			})
		}
	
		customerResponse.Image = url
	}
	ctx.JSON(http.StatusCreated, customerResponse)
}

func (c *Controller) ListCustomersHandle(ctx *gin.Context) {
	customers, err := c.customerUsecase.ListCustomers()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customers)
}

func (c *Controller) GetUserByFieldHandle(ctx *gin.Context) {
	var fieldName = ctx.Param("fieldName")
	var value = ctx.Param("value")

	customer, err := c.customerUsecase.GetCustomerByField(fieldName, value)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, customer)
}
