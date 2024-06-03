package helpers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secret_key") // Replace with your secret key

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat": time.Now().Unix(),
		"uid": userID,
	})
	return claims.SignedString(mySigningKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	fmt.Println("Received token string:", tokenString)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println("Unexpected signing method:", token.Method)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		fmt.Println("Token parsed successfully")
		return mySigningKey, nil
	})
}
