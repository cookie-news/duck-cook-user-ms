package usecase

import (
	"duck-cook-user-ms/api/repository"
	"duck-cook-user-ms/entity"
	"duck-cook-user-ms/repository/auth_repository"
	"errors"
	"mime/multipart"
	"regexp"
	"strings"
)

type CustomerUsecase interface {
	ListCustomers() ([]entity.CustomerResponse, error)
	GetCustomerByField(fieldName string, value string) (entity.CustomerResponse, error)
	CreateCustomer(customer entity.Customer) (entity.CustomerResponse, error)
	Validate(customer entity.Customer) (err error)
	UploadImage(image multipart.FileHeader, user string) (string, error)
	DeleteCustomer(idCustomer string) error
}

type customerUsecaseImpl struct {
	customerRepository repository.CustomerRepository
	customerStorage    repository.CustomerStorage
	authRepository     auth_repository.AuthRepository
}

func (usecase customerUsecaseImpl) DeleteCustomer(idCustomer string) error {
	return usecase.customerRepository.DeleteCustomer(idCustomer)
}

func (usecase customerUsecaseImpl) ListCustomers() (customersList []entity.CustomerResponse, err error) {
	result, err := usecase.customerRepository.ListCustomers()
	var customers []entity.CustomerResponse
	if err != nil {
		return
	}

	for _, customer := range result {
		imageUrl := usecase.customerStorage.GetPublicUrl(customer.User)
		customer.Image = imageUrl
		customers = append(customers, customer)
	}

	return customers, err
}

func (usecase customerUsecaseImpl) CreateCustomer(customer entity.Customer) (customerResult entity.CustomerResponse, err error) {
	err = usecase.Validate(customer)

	if err != nil {
		return
	}

	err = usecase.authRepository.CreateUser(customer.User, customer.Pass, customer.Email)

	if err != nil {
		return
	}

	return usecase.customerRepository.CreateCustomer(customer)
}

func (usecase customerUsecaseImpl) GetCustomerByField(fieldName string, value string) (customer entity.CustomerResponse, err error) {
	if fieldName == "" || value == "" {
		return customer, err
	}

	result, err := usecase.customerRepository.GetCustomerByField(fieldName, value)
	if err != nil {
		return
	}

	imageUrl := usecase.customerStorage.GetPublicUrl(result.User)
	result.Image = imageUrl
	return result, err
}

func (usecase customerUsecaseImpl) UploadImage(image multipart.FileHeader, user string) (url string, err error) {
	key, err := usecase.customerStorage.UploadImage(image, user)
	if err != nil {
		return
	}
	parts := strings.Split(key, "/")
	filename := parts[len(parts)-1]

	url = usecase.customerStorage.GetPublicUrl(filename)

	return url, err
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

func NewCustomerUseCase(
	customerRepository repository.CustomerRepository,
	customerStorage repository.CustomerStorage,
	authRepository auth_repository.AuthRepository) CustomerUsecase {
	return &customerUsecaseImpl{
		customerRepository,
		customerStorage,
		authRepository,
	}
}
