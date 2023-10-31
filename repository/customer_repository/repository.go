package customerrepository

import (
	"context"
	"duck-cook-user-ms/api/repository"
	"duck-cook-user-ms/db"
	"duck-cook-user-ms/entity"
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

func (r repositoryImpl) GetCustomerByEmail(email string) (customer entity.Customer, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customerModel Customer
	var customerEntity entity.Customer

	filter := bson.M{"email": email}

	err = r.customerCollection.FindOne(ctx, filter).Decode(&customerModel)
	if err != nil {
		return customerEntity, err
	}

	return customerModel.ToEntityCustomer(), err
}

func (r repositoryImpl) GetCustomerByUser(user string) (customer entity.Customer, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customerModel Customer
	var customerEntity entity.Customer

	filter := bson.M{"user": user}

	err = r.customerCollection.FindOne(ctx, filter).Decode(&customerModel)
	if err != nil {
		return customerEntity, err
	}

	return customerModel.ToEntityCustomer(), err
}

func (r repositoryImpl) GetCustomerByField(fieldName string, value string) (customer entity.Customer, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customerModel Customer
	var customerEntity entity.Customer

	filter := bson.M{fieldName: value}

	err = r.customerCollection.FindOne(ctx, filter).Decode(&customerModel)
	if err != nil {
		return customerEntity, err
	}

	return customerModel.ToEntityCustomer(), err
}

func (r repositoryImpl) ListCustomers() (customer []entity.Customer, err error) {
	ctx, cancel := doTimeout()

	defer cancel()

	var customersModel []Customer
	var customerEntity []entity.Customer = []entity.Customer{}

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
		customerEntity = append(customerEntity, customer.ToEntityCustomer())
	}

	return customerEntity, nil
}

func (r repositoryImpl) CreateCustomer(customer entity.Customer) (entity.Customer, error) {
	ctx, cancel := doTimeout()
	defer cancel()

	var customerModel Customer
	customerModel = *customerModel.FromEntity(customer)
	passHash, err := HashPassword(customer.Pass)

	if err != nil {
		return customerModel.ToEntityCustomer(), err
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
							return customerModel.ToEntityCustomer(), errors.New("duplicate " + fieldName)
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

	return customerModel.ToEntityCustomer(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func New(mongoDb mongo.Database) repository.CustomerRepository {
	customerCollection := mongoDb.Collection(db.COLLECTION_CUSTOMER)
	return &repositoryImpl{customerCollection}
}
