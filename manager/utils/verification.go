package utils

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/b1izko/test-pizza/manager/store"
	"github.com/b1izko/test-pizza/manager/store/admin"
)

// AuthToken verifies the user by JWT in HTTP header
func AuthToken(r *http.Request, s store.Storage) (*admin.Model, error) {
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		log.Printf("%v", token)
		return nil, errors.New("Invalid token type")
	}

	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := admin.GetClaims(token, s)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("Invalid token")
	}

	if err := claims.Valid(); err != nil {
		log.Printf("%v", err)
		return nil, errors.New("Invalid token")
	}

	admin, err := admin.ByID(claims.ID, s)
	if err != nil {
		log.Printf("%v", err)
		return nil, errors.New("Invalid user")
	}

	return admin, nil
}
