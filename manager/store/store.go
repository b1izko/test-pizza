package store

import (
	"context"
	"time"

	"github.com/b1izko/test-pizza/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage interface that describes application store
type Storage interface {
	Database() *mongo.Database
	Connect() error
	Disconnect() error
}

// Store is a struct that controlls database operations
type Store struct {
	URL    string
	DBName string
	client *mongo.Client
}

// New creates a new storage
func New(URL, Username, Password, DBName string) (*Store, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URL))
	//TODO
	//client, err := mongo.NewClient(options.Client().ApplyURI(URL).SetAuth(options.Credential{
	//	Username: Username,
	//	Password: Password,
	//}))

	if logger.IsError(err, "Failed to create a new storage") {
		return nil, err
	}

	s := &Store{
		URL:    URL,
		DBName: DBName,
		client: client,
	}

	return s, nil
}

// Database returns current database
func (s *Store) Database() *mongo.Database {
	return s.client.Database(s.DBName)
}

// Connect to the database
func (s *Store) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := s.client.Connect(ctx)
	if logger.IsError(err, "Failed to create a new storage") {
		return err
	}
	err = s.client.Ping(context.TODO(), nil)
	if logger.IsError(err, "Failed to create a new storage") {
		return err
	}

	return nil
}

// Disconnect store
func (s *Store) Disconnect() error {
	return s.client.Disconnect(context.TODO())
}
