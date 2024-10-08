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

func (s *CustomerService) AddMagicEmailToken(
	ctx context.Context, magicToken models.EmailMagicToken) (models.EmailMagicToken, error) {
	identity, err := s.repo.InsertEmailMagicToken(ctx, magicToken)
	if err != nil {
		return models.EmailMagicToken{}, ErrCustomer
	}
	return identity, nil
}

func (s *CustomerService) RemoveMagicEmailToken(
	ctx context.Context, workspaceId string, customerId string, email string) error {
	err := s.repo.DeleteMagicEmailToken(ctx, workspaceId, customerId, email)
	if err != nil {
		return ErrCustomer
	}
	return nil
}

func (s *CustomerService) HasProvidedEmailIdentity(ctx context.Context, workspaceId string, customerId string) (bool, error) {
	exists, err := s.repo.EmailIdentityExists(ctx, workspaceId, customerId)
	if err != nil {
		return false, ErrEmailIdentityCheck
	}
	return exists, nil
}

func (s *CustomerService) GenerateKycMailVerifyToken(
	sk string, workspaceId string, customerId string, email string, expiresAt time.Time, redirectUrl string) (string, error) {
	claims := models.KycMailJWTClaims{
		WorkspaceId: workspaceId,
		Email:       email,
		RedirectUrl: redirectUrl,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "kyc.zyg.ai",
			Subject:   customerId,
			Audience:  []string{"kyc"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ID:        workspaceId + ":" + customerId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	j, err := token.SignedString([]byte(sk))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token got error: %v", err)
	}
	return j, nil
}

func (s *CustomerService) VerifyKycMailToken(token string, hmacSecret []byte) (models.KycMailJWTClaims, error) {
	t, err := jwt.ParseWithClaims(
		token, &models.KycMailJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%v", token.Header["alg"])
			}
			return hmacSecret, nil
		})

	if err != nil {
		return models.KycMailJWTClaims{}, fmt.Errorf("%v", err)
	} else if claims, ok := t.Claims.(*models.KycMailJWTClaims); ok {
		return *claims, nil
	}
	return models.KycMailJWTClaims{}, fmt.Errorf("error parsing jwt token")
}

func (s *CustomerService) GetKycMailToken(ctx context.Context, token string) (models.EmailMagicToken, error) {
	magicToken, err := s.repo.LookupEmailMagicTokenByToken(ctx, token)
	if err != nil {
		return models.EmailMagicToken{}, ErrKycMailToken
	}
	now := time.Now().UTC()
	// check if the current time is after the token expiration time.
	if now.After(magicToken.ExpiresAt) {
		return models.EmailMagicToken{}, ErrKycMailTokenExpired
	}

	return magicToken, err
}
