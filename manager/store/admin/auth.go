package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/b1izko/test-pizza/manager/store"

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "+f7FX3):FoBNc28Wh81c1-3z8STpQLGcFiK8I-z0!FzvtPwi9tbwid*UJblB;V8c(c*P_HjoCW)sDtF%o?L=EmIkd+_Bk6l(I9FAIu1_S!"

// JWToken ...
type JWToken struct {
	jwt.StandardClaims
	ID       string
	Status   int
	LastAuth time.Time
}

// CheckToken ...
func CheckToken(rawToken string) (*JWToken, error) {
	if !strings.HasPrefix(rawToken, "Bearer ") {
		return nil, errors.New("invalid token type")
	}
	tokenString := strings.TrimPrefix(rawToken, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &JWToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWToken); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token expired")
}

// GetClaims decode and verify JWT token
func GetClaims(tokenString string, store store.Storage) (*JWToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWToken); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}
