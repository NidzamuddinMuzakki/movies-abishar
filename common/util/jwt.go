package util

import (
	"fmt"
	"time"

	"github.com/NidzamuddinMuzakki/movies-abishar/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Token struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Exp      string
}

func GenerateTokenPair(name string) (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS512)
	id := uuid.New()
	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["username"] = name
	claims["uuid"] = id.String()
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(config.Cold.JwtSecretKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["username"] = name
	rtClaims["uuid"] = id.String()
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// rt, err := refreshToken.SignedString([]byte(config.Cold.JwtSecretKey))
	// if err != nil {
	// 	return nil, err
	// }

	return map[string]string{
		"token": t,
		// "refreshtoken": rt,
	}, nil
}

func ReadDataToken(token string) Token {
	claims := jwt.MapClaims{}
	var tokenReturn Token
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cold.JwtSecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for key, val := range claims {
		if key == "exp" {
			tokenReturn.Exp = fmt.Sprintf("%v", val)
		}
		if key == "username" {
			tokenReturn.Username = fmt.Sprintf("%v", val)
		}
		if key == "uuid" {
			tokenReturn.Uuid = fmt.Sprintf("%v", val)
		}

	}
	return tokenReturn
}
