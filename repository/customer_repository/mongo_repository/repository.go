package mongodb_repository

import (
	"context"
	"duck-cook-user-ms/api/repository"
	"duck-cook-user-ms/entity"
	"duck-cook-user-ms/pkg/mongodb"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type repositoryImpl struct {
	customerCollection *mongo.Collection
}

func doTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (r repositoryImpl) GetCustomerByField(fieldName string, value string) (customer entity.CustomerResponse, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customerModel Customer

	filter := bson.M{fieldName: value}

	err = r.customerCollection.FindOne(ctx, filter).Decode(&customerModel)
	if err != nil {
		return entity.CustomerResponse{}, err
	}

	return customerModel.ToEntityCustomerResponse(), err
}

func (r repositoryImpl) ListCustomers() (customer []entity.CustomerResponse, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customersModel []Customer
	var customerEntity []entity.CustomerResponse = []entity.CustomerResponse{}

	cursor, err := r.customerCollection.Find(ctx, bson.D{})
	if err != nil {
		if err == mongo.ErrNilDocument {
			return customerEntity, nil
		}
		return nil, err
	}

	if err = cursor.All(context.TODO(), &customersModel); err != nil {
		return nil, err
	}

	for _, customer := range customersModel {
		customerEntity = append(customerEntity, customer.ToEntityCustomerResponse())
	}

	return customerEntity, nil
}

func (r repositoryImpl) CreateCustomer(customer entity.Customer) (entity.CustomerResponse, error) {
	ctx, cancel := doTimeout()
	defer cancel()

	var customerModel Customer
	customerModel = *customerModel.FromEntity(customer)
	passHash, err := HashPassword(customer.Pass)

	if err != nil {
		return customerModel.ToEntityCustomerResponse(), err
	}

	customerModel.Pass = passHash
	res, err := r.customerCollection.InsertOne(ctx, &customerModel)

	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, writeErr := range writeErr.WriteErrors {
				if writeErr.Code == 11000 {
					errorMsg := writeErr.Message
					startIdx := strings.Index(errorMsg, "{")
					endIdx := strings.Index(errorMsg, "}")
					if startIdx != -1 && endIdx != -1 {
						fieldInfo := errorMsg[startIdx+1 : endIdx]

						re := regexp.MustCompile(`(\w+):`)
						match := re.FindStringSubmatch(fieldInfo)
						if len(match) >= 2 {
							fieldName := match[1]
							return customerModel.ToEntityCustomerResponse(), errors.New("duplicate " + fieldName)
						}
					}

				} else {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
	}

	customerModel.ID = res.InsertedID.(primitive.ObjectID)

	return customerModel.ToEntityCustomerResponse(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func New(mongoDb mongo.Database) repository.CustomerRepository {
	customerCollection := mongoDb.Collection(mongodb.COLLECTION_CUSTOMER)
	return &repositoryImpl{customerCollection}
}
