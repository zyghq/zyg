package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zyghq/zyg"
	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
	"github.com/zyghq/zyg/services/tasks"
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
		WorkspaceId:     c.WorkspaceId,
		ExternalId:      externalId,
		Email:           email,
		Phone:           phone,
		IsEmailVerified: c.IsEmailVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth.zyg.ai",
			Subject:   c.CustomerId,
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().AddDate(1, 0, 0)), // 1 year from now
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

func (s *CustomerService) AddClaimedMail(
	ctx context.Context, claimed models.ClaimedMail) (models.ClaimedMail, error) {
	claim, err := s.repo.InsertClaimedMail(ctx, claimed)
	if err != nil {
		return models.ClaimedMail{}, ErrClaimedMail
	}
	return claim, nil
}

func (s *CustomerService) RemoveCustomerClaimedMail(
	ctx context.Context, workspaceId string, customerId string, email string) error {
	err := s.repo.DeleteCustomerClaimedMail(ctx, workspaceId, customerId, email)
	if err != nil {
		return ErrClaimedMail
	}
	return nil
}

func (s *CustomerService) GetRecentValidClaimedMail(
	ctx context.Context, workspaceId string, customerId string) (string, error) {
	claimed, err := s.repo.LookupLatestClaimedMail(ctx, workspaceId, customerId)
	if errors.Is(err, repository.ErrEmpty) {
		return "", ErrClaimedMailNotFound
	}
	if err != nil {
		return "", ErrClaimedMail
	}
	now := time.Now().UTC()
	// check if the current time is after the token expiration time.
	if now.After(claimed.ExpiresAt) {
		return "", ErrClaimedMailExpired
	}
	return claimed.Email, nil
}

func (s *CustomerService) GenerateMailVerificationToken(
	sk string, workspaceId string, customerId string, email string, expiresAt time.Time, redirectUrl string,
) (string, error) {
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

func (s *CustomerService) VerifyMailVerificationToken(hmacSecret []byte, token string,
) (models.KycMailJWTClaims, error) {
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

func (s *CustomerService) GetValidClaimedMailByToken(
	ctx context.Context, token string) (models.ClaimedMail, error) {
	claimed, err := s.repo.LookupClaimedMailByToken(ctx, token)

	if errors.Is(err, repository.ErrEmpty) {
		return models.ClaimedMail{}, ErrClaimedMailNotFound
	}
	if err != nil {
		return models.ClaimedMail{}, ErrClaimedMail
	}
	now := time.Now().UTC()
	// check if the current time is after the token expiration time.
	if now.After(claimed.ExpiresAt) {
		return models.ClaimedMail{}, ErrClaimedMailExpired
	}
	return claimed, err
}

func (s *CustomerService) ClaimMailForVerification(
	ctx context.Context, customer models.Customer, sk string,
	email string, name *string, hasConflict bool, contextMessage string, redirectTo string,
) (models.ClaimedMail, error) {
	expiresAt := time.Now().UTC().AddDate(0, 0, 2) // 2 days
	jt, err := s.GenerateMailVerificationToken(
		sk, customer.WorkspaceId, customer.CustomerId,
		email, expiresAt, redirectTo,
	)
	if err != nil {
		return models.ClaimedMail{}, ErrClaimedMail
	}

	claim := models.ClaimedMail{}.NewVerification(
		customer.WorkspaceId, customer.CustomerId,
		email, hasConflict, expiresAt, jt,
	)
	claim, err = s.AddClaimedMail(ctx, claim)
	if err != nil {
		return models.ClaimedMail{}, err
	}

	if customer.IsVisitor() {
		dup := customer.MakeCopy()
		dup.Role = dup.Lead()
		if name != nil {
			dup.Name = *name
		}
		dup, err = s.UpdateCustomer(ctx, dup)
		if err != nil {
			return models.ClaimedMail{}, ErrCustomer
		}
	}
	verifyLink := zyg.GetXServerUrl() + "/mail/kyc/?t=" + claim.Token
	err = tasks.SendKycMail(claim.Email, contextMessage, verifyLink)
	slog.Error("tasks send kyc mail failed", slog.Any("err", err))
	return claim, nil
}

func (s *CustomerService) AppendEvent(
	ctx context.Context, event models.CustomerEvent) (models.CustomerEvent, error) {
	event, err := s.repo.InsertEvent(ctx, event)
	if err != nil {
		return models.CustomerEvent{}, ErrCustomerEvent
	}
	return event, nil
}

func (s *CustomerService) ListEvents(ctx context.Context, customerId string) ([]models.CustomerEvent, error) {
	events, err := s.repo.FetchEventsByCustomerId(ctx, customerId)
	if err != nil {
		return []models.CustomerEvent{}, ErrCustomerEvent
	}
	return events, nil
}
