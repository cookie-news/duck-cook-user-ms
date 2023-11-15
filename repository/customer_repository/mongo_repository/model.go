package mongodb_repository

import (
	"duck-cook-user-ms/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
	Email     string             `bson:"email"`
	User      string             `bson:"user"`
	Name      string             `bson:"name"`
}

// Transforma o modelo do mongo para uma entidade
func (customer Customer) ToEntityCustomerResponse() entity.CustomerResponse {
	return entity.CustomerResponse{
		ID:    customer.ID.Hex(),
		Email: customer.Email,
		User:  customer.User,
		Name:  customer.Name,
	}
}

// Transforma da entidade para o modelo do mongo
func (Customer) FromEntity(customer entity.Customer) *Customer {
	id, _ := primitive.ObjectIDFromHex(customer.ID)
	return &Customer{
		ID:    id,
		Email: customer.Email,
		User:  customer.User,
		Name:  customer.Name,
	}
}
