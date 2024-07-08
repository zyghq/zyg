package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type CustomerService struct {
	repo ports.CustomerRepositorer
}

func NewCustomerService(repo ports.CustomerRepositorer) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) GetCustomerByExternalId(ctx context.Context, workspaceId string, externalId string,
) (models.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByExtId(ctx, workspaceId, externalId)

	if errors.Is(err, repository.ErrEmpty) {
		return customer, ErrCustomerNotFound
	}

	if err != nil {
		return customer, ErrCustomer
	}
	return customer, nil
}

func (s *CustomerService) GetCustomerByEmail(ctx context.Context, workspaceId string, email string,
) (models.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByEmail(ctx, workspaceId, email)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) GetCustomerByPhone(ctx context.Context, workspaceId string, email string,
) (models.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByPhone(ctx, workspaceId, email)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) GenerateCustomerToken(c models.Customer) (string, error) {
	var externalId string
	var email string
	var phone string

	audience := []string{"customer"}

	sk, err := zyg.GetEnv("ZYG_CUSTOMER_JWT_SECRET")
	if err != nil {
		return "", fmt.Errorf("failed to get env ZYG_CUSTOMER_JWT_SECRET got error: %v", err)
	}

	if !c.ExternalId.Valid {
		externalId = ""
	} else {
		externalId = c.ExternalId.String
	}

	if !c.Email.Valid {
		email = ""
	} else {
		email = c.Email.String
	}

	if !c.Phone.Valid {
		phone = ""
	} else {
		phone = c.Phone.String
	}

	claims := models.CustomerJWTClaims{
		WorkspaceId: c.WorkspaceId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.zyg.ai",
			Subject:   c.CustomerId,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(1, 0, 0)), // Expires 1 year from now
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        c.WorkspaceId + ":" + c.CustomerId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token got error: %v", err)
	}
	return jwt, nil
}
