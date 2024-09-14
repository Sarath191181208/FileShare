package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
  Id string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(jwtKey, email, id string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
    Id: id,
    Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(jwtKey, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
