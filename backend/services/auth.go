package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

const DefaultAuthProvider string = "supabase"

func ParseJWTToken(token string, hmacSecret []byte) (ac models.AuthJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &models.AuthJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return ac, fmt.Errorf("%v", err)
	} else if claims, ok := t.Claims.(*models.AuthJWTClaims); ok {
		return *claims, nil
	}
	return ac, fmt.Errorf("error parsing jwt token")
}

func ParseCustomerJWTToken(token string, hmacSecret []byte) (cc models.CustomerJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(token, &models.CustomerJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return cc, fmt.Errorf("%v", err)
	} else if claims, ok := t.Claims.(*models.CustomerJWTClaims); ok {
		return *claims, nil
	}
	return cc, fmt.Errorf("error parsing jwt token")
}

type AuthService struct {
	accountRepo ports.AccountRepositorer
}

func NewAuthService(accountRepo ports.AccountRepositorer) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
	}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, authUserId string) (models.Account, error) {
	account, err := s.accountRepo.FetchAccountByAuthId(ctx, authUserId)

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

func (s *AuthService) ValidatePersonalAccessToken(ctx context.Context, token string) (models.Account, error) {
	account, err := s.accountRepo.LookupAccountByToken(ctx, token)

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

type CustomerAuthService struct {
	customerRepo ports.CustomerRepositorer
}

func NewCustomerAuthService(customerRepo ports.CustomerRepositorer) *CustomerAuthService {
	return &CustomerAuthService{
		customerRepo: customerRepo,
	}
}

func (s *CustomerAuthService) GetWorkspaceCustomerIgnoreRole(
	ctx context.Context, workspaceId string, customerId string) (models.Customer, error) {
	customer, err := s.customerRepo.LookupWorkspaceCustomerWithoutRoleById(ctx, workspaceId, customerId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrCustomerNotFound
	}

	if err != nil {
		return models.Customer{}, ErrCustomer
	}
	return customer, nil
}

func (s *CustomerAuthService) GetWidgetLinkedSecretKey(
	ctx context.Context, widgetId string) (models.SecretKey, error) {
	sk, err := s.customerRepo.LookupSecretKeyByWidgetId(ctx, widgetId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.SecretKey{}, ErrSecretKeyNotFound
	}

	if err != nil {
		return models.SecretKey{}, ErrSecretKey
	}
	return sk, nil
}
