package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ReadJson(r *http.Request, input interface{}) error {

	log.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(input)

	if err != nil {
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, input interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(input)
	if err != nil {
		return err
	}

	return nil
}

var secretKey = []byte("secret-key")

func CreateToken(id int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func VerifyToken(tokenString string) (int32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	idFloat, ok := claims["id"].(float64) // JWT stores numbers as float64
	if !ok {
		return 0, fmt.Errorf("invalid token id")
	}

	return int32(idFloat), nil
}
