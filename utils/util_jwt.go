package utils

import (
	"EX_okexquant/mylog"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var hmacSampleSecret []byte

func init() {
	hmacSampleSecret = []byte("bitway-todo_block")
}

func CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		mylog.Logger.Error().Msgf("[CreateToken] token SignedString failed, err:%v, tokenString:%v", err, tokenString)
		return "", err
	}

	return tokenString, err
}

func VerifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		mylog.Logger.Error().Msgf("[VerifyToken] jwt Parse failed, err:%v, token:%v", err, token)
		return false
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false
	}

	return true
}
