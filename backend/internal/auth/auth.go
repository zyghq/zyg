package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg/internal/domain"
)

const DefaultAuthProvider string = "supabase"

func ParseJWTToken(token string, hmacSecret []byte) (ac domain.AuthJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &domain.AuthJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return ac, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*domain.AuthJWTClaims); ok {
		return *claims, nil
	}

	return ac, fmt.Errorf("error parsing jwt token: %v", token)
}

func ParseCustomerJWTToken(token string, hmacSecret []byte) (cc domain.CustomerJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &domain.CustomerJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return cc, fmt.Errorf("error validating jwt token with error: %v", err)
	} else if claims, ok := t.Claims.(*domain.CustomerJWTClaims); ok {
		return *claims, nil
	}
	return cc, fmt.Errorf("error parsing jwt token: %v", token)
}
