package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg/internal/model"
)

const DefaultAuthProvider string = "supabase"

func ParseJWTToken(token string, hmacSecret []byte) (ac model.AuthJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &model.AuthJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return ac, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*model.AuthJWTClaims); ok {
		return *claims, nil
	}

	return ac, fmt.Errorf("error parsing jwt token: %v", token)
}

func ParseCustomerJWTToken(token string, hmacSecret []byte) (cc model.CustomerJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &model.CustomerJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return cc, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*model.CustomerJWTClaims); ok {
		return *claims, nil
	}
	return cc, fmt.Errorf("error parsing jwt token: %v", token)
}
