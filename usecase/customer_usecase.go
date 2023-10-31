package usecase

import (
	"duck-cook-user-ms/api/repository"
	"duck-cook-user-ms/entity"
	"errors"
	"regexp"
)

type CustomerUsecase interface {
	ListCustomers() ([]entity.Customer, error)
	GetCustomerByField(fieldName string, value string) (entity.Customer, error)
	CreateCustomer(customer entity.Customer) (entity.Customer, error)
	Validate(customer entity.Customer) (err error)
}

type customerUsecaseImpl struct {
	customerRepository repository.CustomerRepository
}

func (usecase customerUsecaseImpl) ListCustomers() (customersList []entity.Customer, err error) {
	customers, err := usecase.customerRepository.ListCustomers()
	if err != nil {
		return
	}
	return customers, err
}

func (usecase customerUsecaseImpl) CreateCustomer(customer entity.Customer) (customerResult entity.Customer, err error) {
	err = usecase.Validate(customer)

	if err != nil {
		return
	}

	return usecase.customerRepository.CreateCustomer(customer)
}

func (usecase customerUsecaseImpl) GetCustomerByField(fieldName string, value string) (customer entity.Customer, err error) {

	result, err := usecase.customerRepository.GetCustomerByField(fieldName, value)
	if err != nil {
		return
	}

	return result, err
}

var (
	ErrInvalidEmail       = errors.New("email is invalid")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUserAlreadyExists  = errors.New("user alrady registered")
)

func (customerUsecaseImpl *customerUsecaseImpl) Validate(customer entity.Customer) (err error) {
	emailIsValid := checkEmail(customer.Email)
	emailAlreadyExists := customerUsecaseImpl.emailExists(customer.Email)
	userAlreadyExists := customerUsecaseImpl.userExists(customer.User)

	if emailAlreadyExists {
		return ErrEmailAlreadyExists
	}

	if userAlreadyExists {
		return ErrUserAlreadyExists
	}

	if !emailIsValid {
		return ErrInvalidEmail
	}

	return nil
}

func (customerUsecase customerUsecaseImpl) userExists(user string) bool {
	result, _ := customerUsecase.GetCustomerByField("user", user)
	return result.User != ""
}

func (customerUsecase customerUsecaseImpl) emailExists(email string) bool {
	result, _ := customerUsecase.GetCustomerByField("email", email)
	return result.Email != ""
}

func checkEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func NewCustomerUseCase(customerRepository repository.CustomerRepository) CustomerUsecase {
	return &customerUsecaseImpl{
		customerRepository,
	}
}
