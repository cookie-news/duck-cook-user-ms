package entity

import "mime/multipart"

type Customer struct {
	ID    string                `form:"id"`
	Email string                `form:"email" binding:"required"`
	User  string                `form:"user" binding:"required"`
	Name  string                `form:"name" binding:"required"`
	Pass  string                `form:"pass" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type CustomerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	User  string `json:"user"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (customer Customer) ToResponseEntity() CustomerResponse {
	return CustomerResponse{
		ID:    customer.ID,
		Email: customer.Email,
		User:  customer.User,
		Name:  customer.Name,
		Image: "",
	}
}
