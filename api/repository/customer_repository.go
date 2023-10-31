package repository

import "duck-cook-user-ms/entity"

type CustomerRepository interface {
	CreateCustomer(customer entity.Customer) (entity.Customer, error)
	ListCustomers() ([]entity.Customer, error)
	GetCustomerByEmail(email string) (entity.Customer, error)
	GetCustomerByUser(user string) (entity.Customer, error)
	GetCustomerByField(fieldName string, value string) (entity.Customer, error)
}
