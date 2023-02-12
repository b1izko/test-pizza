package admin

import (
	"context"
	"encoding/json"
	"strings"
	"time"

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
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "id" {
			m.ID, err = primitive.ObjectIDFromHex(v.(string))
			if err != nil {
				return err
			}
		}

		if strings.ToLower(k) == "status" {
			m.Status = int(v.(float64))
		}

		if strings.ToLower(k) == "last_auth" {
			t, err := time.Parse(TimeTemplate, v.(string))
			if err != nil {
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
	if err != nil {
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
	if err != nil {
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
		return errors.New("Invalid ID")
	}

	if m.Login == "root" {
		return errors.New("cannot remove root user")
	}

	result, err := store.Database().Collection(CollectionName).DeleteOne(ctx, bson.M{"_id": m.ID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("Delete error")
	}

	return nil
}

//ByID returns admin by ID
func ByID(id string, store store.Storage) (*Model, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.WithMessage(err, "cannot parse ID")
	}

	var usr Model
	if err := store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"_id": _id}).Decode(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

//ByLogin returns admin by login
func ByLogin(login string, store store.Storage) (*Model, error) {
	var usr Model
	if err := store.Database().Collection(CollectionName).FindOne(ctx, bson.M{"login": login}).Decode(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

//Find returns all model by filter
func Find(store store.Storage, fltr *Model, offset, limit int64) ([]*Model, error) {

	filter, err := makeFilter(fltr)
	if err != nil {
		return nil, err
	}

	cur, err := store.Database().Collection(CollectionName).Find(ctx, *filter)
	if err != nil {
		return nil, err
	}

	auth := []*Model{}

	if err := cur.All(ctx, &auth); err != nil {
		return nil, err
	}

	return auth, nil
}

// makeFilter returns filter for Collection.Find()
func makeFilter(values *Model) (*bson.D, error) {
	data, err := values.MarshalJSON()
	if err != nil {
		return nil, err
	}

	filter := &bson.D{}
	err = bson.Unmarshal(data, filter)
	if err != nil {
		return nil, err
	}

	return filter, nil

}

// Check return model for login:password
func Check(login string, password string, store store.Storage) (*Model, error) {
	model, err := ByLogin(login, store)
	if err != nil {
		return nil, errors.New("Wrong login")
	}

	if model.Password != password {
		return nil, errors.New("Wrong password")
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

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// MakeFilter returns filter for Collection.Find()
func MakeFilter(values *map[string]string) *bson.D {
	filter := bson.D{}

	var keys []string

	for key := range *values {
		keys = append(keys, key)
	}

	for _, key := range keys {
		if (*values)[key] != "" {
			if key != "offset" && key != "limit" {
				filter = append(filter, bson.E{Key: key, Value: (*values)[key]})
			}
		}
	}
	return &filter
}
