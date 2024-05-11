package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
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

type AuthService struct {
	repo ports.AccountRepositorer
}

func NewAuthService(repo ports.AccountRepositorer) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) GetAuthUser(ctx context.Context, authUserId string) (domain.Account, error) {
	account, err := s.repo.GetByAuthUserId(ctx, authUserId)

	if errors.Is(err, repository.ErrQuery) {
		return account, ErrAccount
	}

	if errors.Is(err, repository.ErrEmpty) {
		return account, ErrAccountNotFound
	}

	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *AuthService) GetPatAccount(ctx context.Context, token string) (domain.Account, error) {
	account, err := s.repo.GetAccountByToken(ctx, token)

	if errors.Is(err, repository.ErrQuery) {
		return account, ErrAccount
	}

	if errors.Is(err, repository.ErrEmpty) {
		return account, ErrAccountNotFound
	}

	if err != nil {
		return account, err
	}
	return account, nil
}
