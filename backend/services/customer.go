package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type CustomerService struct {
	repo ports.CustomerRepositorer
}

func NewCustomerService(
	repo ports.CustomerRepositorer) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) GenerateCustomerJwt(c models.Customer, sk string) (string, error) {
	var externalId, email, phone *string
	audience := []string{"customer"}

	if c.ExternalId.Valid {
		externalId = &c.ExternalId.String
	}
	if c.Email.Valid {
		email = &c.Email.String
	}
	if c.Phone.Valid {
		phone = &c.Phone.String
	}

	claims := models.CustomerJWTClaims{
		WorkspaceId: c.WorkspaceId,
		ExternalId:  externalId,
		Email:       email,
		Phone:       phone,
		IsVerified:  c.IsVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.zyg.ai",
			Subject:   c.CustomerId,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().AddDate(1, 0, 0)), // Expires 1 year from now
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ID:        c.WorkspaceId + ":" + c.CustomerId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	j, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token got error: %v", err)
	}
	return j, nil
}

func (s *CustomerService) VerifyExternalId(sk string, hash string, externalId string) bool {
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(externalId))
	hashHex := hex.EncodeToString(h.Sum(nil))
	return hashHex == hash
}

func (s *CustomerService) VerifyEmail(sk string, hash string, email string) bool {
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(email))
	hashHex := hex.EncodeToString(h.Sum(nil))
	return hashHex == hash
}

func (s *CustomerService) VerifyPhone(sk string, hash string, phone string) bool {
	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(phone))
	hashHex := hex.EncodeToString(h.Sum(nil))
	return hashHex == hash
}

func (s *CustomerService) UpdateCustomer(
	ctx context.Context, customer models.Customer) (models.Customer, error) {
	customer, err := s.repo.ModifyCustomerById(ctx, customer)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrCustomerNotFound
	}
	if err != nil {
		return models.Customer{}, ErrCustomer
	}
	return customer, nil
}

func (s *CustomerService) AddCustomerEmailIdentity(
	ctx context.Context, emailIdentity models.EmailIdentity) (models.EmailIdentity, error) {
	identity, err := s.repo.InsertEmailIdentity(ctx, emailIdentity)
	if err != nil {
		return models.EmailIdentity{}, ErrCustomer
	}
	return identity, nil
}

func (s *CustomerService) HasProvidedEmailIdentity(ctx context.Context, customerId string) (bool, error) {
	exists, err := s.repo.EmailIdentityExists(ctx, customerId)
	if err != nil {
		return false, ErrEmailIdentityCheck
	}
	return exists, nil
}
