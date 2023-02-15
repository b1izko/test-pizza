package admin

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/b1izko/test-pizza/internal/logger"
	"github.com/b1izko/test-pizza/manager/store"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionName for authorization
const CollectionName = "admins"

var ctx = context.TODO()

// Statuses
const (
	StatusRoot = iota
	StatusAdmin
	StatusModerator
)

// TimeTemplate for Model.Time
const TimeTemplate = time.RFC822 // "2006.01.02 15:04:05"

// Model for authorization
type Model struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login    string             `json:"login" bson:"login"`
	Password string             `json:"password" bson:"password"`
	Status   int                `json:"status" bson:"status"`
	LastAuth time.Time          `json:"last_auth" bson:"last_auth"`
}

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

		if strings.ToLower(k) == "status" {
			m.Status = int(v.(float64))
		}

		if strings.ToLower(k) == "last_auth" {
			t, err := time.Parse(TimeTemplate, v.(string))
			if logger.IsError(err, "Failed to unmarshal JSON") {
				return err
			}
			m.LastAuth = t
		}

	}

	var buffer struct {
		Login    string `json:"login" bson:"login"`
		Password string `json:"password" bson:"password"`
	}

	err = json.Unmarshal(j, &buffer)
	if logger.IsError(err, "Failed to unmarshal JSON") {
		return err
	}
	m.Login = buffer.Login
	m.Password = buffer.Password

	return nil
}

func (m Model) MarshalJSON() ([]byte, error) {
	resume := struct {
		ID       string `json:"id" bson:"_id,omitempty"`
		Login    string `json:"login" bson:"login"`
		Password string `json:"password" bson:"password"`
		Status   int    `json:"status" bson:"status"`
		LastAuth string `json:"last_auth" bson:"last_auth"`
	}{
		ID:       m.ID.Hex(),
		Login:    m.Login,
		Password: m.Password,
		LastAuth: m.LastAuth.Format(TimeTemplate),
		Status:   m.Status,
	}

	return json.Marshal(resume)
}

// Save the model
func (m *Model) Save(store store.Storage) error {
	f := bson.M{"login": m.Login}
	if !m.ID.IsZero() {
		f = bson.M{"_id": m.ID}
	}

	result, err := store.Database().Collection(CollectionName).UpdateOne(ctx, f, bson.M{"$set": m}, options.Update().SetUpsert(true))
	if logger.IsError(err, "Failed to save a user") {
		return err
	}

	if result.UpsertedCount > 0 {
		m.ID = result.UpsertedID.(primitive.ObjectID)
	}

	return nil
}

// Remove the model
func (m *Model) Remove(store store.Storage) error {
	if m.ID.IsZero() {
		err := errors.New("User not found")
		logger.IsError(err, "Failed to remove a user")
		return err
	}

	if m.Login == "root" {
		err := errors.New("cannot remove root user")
		logger.IsError(err, "Failed to remove a user")
		return err
	}

	result, err := store.Database().Collection(CollectionName).DeleteOne(ctx, bson.M{"_id": m.ID})
	if logger.IsError(err, "Failed to remove a user") {
		return err
	}

	if result.DeletedCount == 0 {
		err := errors.New("User not deleted")
		logger.IsError(err, "Failed to remove user")
		return err
	}

	return nil
}

// ByID returns admin by ID
func ByID(id string, store store.Storage) (*Model, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if logger.IsError(err, "Failed to parse the user by ID") {
		return nil, err
	}

	var user Model
	err = store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"_id": _id}).Decode(&user)
	if logger.IsError(err, "Failed to parse the user by ID") {
		return nil, err
	}

	return &user, nil
}

// ByLogin returns admin by login
func ByLogin(login string, store store.Storage) (*Model, error) {
	var user Model
	err := store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if logger.IsError(err, "Failed to parse the user by login") {
		return nil, err
	}

	return &user, nil
}

// Check return model for login:password
func Check(login string, password string, store store.Storage) (*Model, error) {
	model, err := ByLogin(login, store)
	if logger.IsError(err, "Wrong login/password") {
		return nil, err
	}

	if logger.IsError(err, "Wrong login/password") {
		return nil, err
	}

	return model, nil
}

// GetJWT returns JWT token
func (m *Model) GetJWT() (string, error) {
	claims := JWToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
		ID:     m.ID.Hex(),
		Status: m.Status,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	if logger.IsError(err, "Getting a JWT ended unsuccessfully") {
		return "", err
	}

	return tokenString, nil
}
