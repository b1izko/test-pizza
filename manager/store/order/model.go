package order

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/b1izko/test-pizza/internal/logger"
	"github.com/b1izko/test-pizza/manager/store"
	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionName for orders
const CollectionName = "orders"

var ctx = context.TODO()

// Model for order
type Model struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User    string             `json:"user,omitempty" bson:"user,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Size    string             `json:"size,omitempty" bson:"size,omitempty"`
	Dough   string             `json:"dough,omitempty" bson:"dough,omitempty"`
	Extra   string             `json:"extra,omitempty" bson:"extra,omitempty"`
	Address string             `json:"address,omitempty" bson:"address,omitempty"`
	Comment string             `json:"comment,omitempty" bson:"comment,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler interface for Model
func (m *Model) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]interface{}
	err := json.Unmarshal(j, &rawStrings)
	if logger.IsError(err, "Failed to unmarshal JSON") {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			m.ID, err = primitive.ObjectIDFromHex(v.(string))
			if logger.IsError(err, "Failed to unmarshal JSON") {
				return err
			}
		}
	}

	var buffer struct {
		User    string `json:"user,omitempty" bson:"user,omitempty"`
		Name    string `json:"name,omitempty" bson:"name,omitempty"`
		Size    string `json:"size,omitempty" bson:"size,omitempty"`
		Dough   string `json:"dough,omitempty" bson:"dough,omitempty"`
		Extra   string `json:"extra,omitempty" bson:"extra,omitempty"`
		Address string `json:"address,omitempty" bson:"address,omitempty"`
		Comment string `json:"comment,omitempty" bson:"comment,omitempty"`
	}

	err = json.Unmarshal(j, &buffer)
	if logger.IsError(err, "Failed to unmarshal JSON") {
		return err
	}
	m.User = buffer.User
	m.Name = buffer.Name
	m.Size = buffer.Size
	m.Dough = buffer.Dough
	m.Extra = buffer.Extra
	m.Address = buffer.Address
	m.Comment = buffer.Comment

	return nil
}

// MarshalJSON initializes custom Marshal
func (m Model) MarshalJSON() ([]byte, error) {
	resume := struct {
		ID      string `json:"id,omitempty" bson:"_id,omitempty"`
		User    string `json:"user,omitempty" bson:"user,omitempty"`
		Name    string `json:"name,omitempty" bson:"name,omitempty"`
		Size    string `json:"size,omitempty" bson:"size,omitempty"`
		Dough   string `json:"dough,omitempty" bson:"dough,omitempty"`
		Extra   string `json:"extra,omitempty" bson:"extra,omitempty"`
		Address string `json:"address,omitempty" bson:"address,omitempty"`
		Comment string `json:"comment,omitempty" bson:"comment,omitempty"`
	}{
		ID:      m.ID.Hex(),
		User:    m.User,
		Name:    m.Name,
		Size:    m.Size,
		Dough:   m.Dough,
		Extra:   m.Extra,
		Address: m.Address,
		Comment: m.Comment,
	}

	return json.Marshal(resume)
}

// Save the order
func (m *Model) Save(store store.Storage) error {
	f := bson.M{"user": m.User}
	if !m.ID.IsZero() {
		f = bson.M{"_id": m.ID}
	}

	result, err := store.Database().Collection(CollectionName).UpdateOne(ctx, f, bson.M{"$set": m}, options.Update().SetUpsert(true))
	if logger.IsError(err, "Failed to save a order") {
		return err
	}

	if result.UpsertedCount > 0 {
		m.ID = result.UpsertedID.(primitive.ObjectID)
	}
	return nil
}

// Remove the order
func (m *Model) Remove(store store.Storage) error {

	if m.ID.IsZero() {
		err := errors.New("Order not found")
		logger.IsError(err, "Failed to remove a order")
		return err
	}

	result, err := store.Database().Collection(CollectionName).DeleteOne(ctx, bson.M{"_id": m.ID})
	if logger.IsError(err, "Failed to remove a order") {
		return err
	}

	if result.DeletedCount == 0 {
		err := errors.New("Order not deleted")
		logger.IsError(err, "Failed to remove order")
		return err
	}

	return nil
}

// ByID returns order by ID
func ByID(id string, store store.Storage) (*Model, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if logger.IsError(err, "Failed to parse the order by ID") {
		return nil, err
	}

	var order Model
	err = store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"_id": _id}).Decode(&order)
	if logger.IsError(err, "Failed to parse the order by ID") {
		return nil, err
	}

	return &order, nil
}

// ByUser returns order by user ID
func ByUser(user string, store store.Storage) (*Model, error) {
	var order Model
	err := store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"user": user}).Decode(&order)
	if logger.IsError(err, "Failed to parse the order by user") {
		return nil, err
	}

	return &order, nil
}
