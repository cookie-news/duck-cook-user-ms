package repository

import "duck-cook-user-ms/entity"

type CustomerRepository interface {
	CreateCustomer(customer entity.Customer) (entity.CustomerResponse, error)
	ListCustomers() ([]entity.CustomerResponse, error)
	GetCustomerByField(fieldName string, value string) (entity.CustomerResponse, error)
}
