package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
)

type CustomerService struct {
	repo ports.CustomerRepositorer
}

func NewCustomerService(repo ports.CustomerRepositorer) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) GetWorkspaceCustomer(ctx context.Context, workspaceId string, customerId string,
) (domain.Customer, error) {
	customer, err := s.repo.GetByWorkspaceCustomerId(ctx, workspaceId, customerId)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) GetWorkspaceCustomerWithExternalId(ctx context.Context, workspaceId string, externalId string,
) (domain.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByExtId(ctx, workspaceId, externalId)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) GetWorkspaceCustomerWithEmail(ctx context.Context, workspaceId string, email string,
) (domain.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByEmail(ctx, workspaceId, email)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) GetWorkspaceCustomerWithPhone(ctx context.Context, workspaceId string, email string,
) (domain.Customer, error) {
	customer, err := s.repo.GetWorkspaceCustomerByPhone(ctx, workspaceId, email)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (s *CustomerService) InitWorkspaceCustomerWithExternalId(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.repo.GetOrCreateCustomerByExtId(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *CustomerService) InitWorkspaceCustomerWithEmail(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.repo.GetOrCreateCustomerByEmail(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *CustomerService) InitWorkspaceCustomerWithPhone(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.repo.GetOrCreateCustomerByPhone(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *CustomerService) IssueJwt(ctx context.Context, c domain.Customer) (string, error) {
	var externalId string
	var email string
	var phone string

	audience := []string{"customer"}

	sk, err := zyg.GetEnv("SUPABASE_JWT_SECRET")
	if err != nil {
		return "", fmt.Errorf("failed to get env SUPABASE_JWT_SECRET got error: %v", err)
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

	claims := domain.CustomerJWTClaims{
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
