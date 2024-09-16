package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateToken(email string, userId primitive.ObjectID) (string ,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"userId": userId,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("AUTH_SECRET")))
} 


func VerifyToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {return nil, errors.New("unexpected auth token signing method")}

		return []byte(os.Getenv("AUTH_SECRET")), nil
	})

	if err != nil {return "", errors.New("could not parse token")}
	if !parsedToken.Valid {return "", errors.New("invalid auth token")}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {return "", errors.New("invalid auth token claims")}

	userId := string(claims["userId"].(string))

	return userId, nil
}