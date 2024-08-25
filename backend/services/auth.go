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

func ParseJWTToken(
	token string, hmacSecret []byte) (ac models.AuthJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(
		token, &models.AuthJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func ParseCustomerJWTToken(
	token string, hmacSecret []byte) (cc models.CustomerJWTClaims, err error) {
	t, err := jwt.ParseWithClaims(
		token, &models.CustomerJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	memberRepo  ports.MemberRepositorer
}

func NewAuthService(
	accountRepo ports.AccountRepositorer, memberRepo ports.MemberRepositorer) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
		memberRepo:  memberRepo,
	}
}

func (s *AuthService) AuthenticateUserAccount(
	ctx context.Context, authUserId string) (models.Account, error) {
	account, err := s.accountRepo.FetchByAuthUserId(ctx, authUserId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Account{}, ErrAccountNotFound
	}

	if err != nil {
		return models.Account{}, ErrAccount
	}

	return account, nil
}

// AuthenticateWorkspaceMember authenticates a workspace member by verifying the existence of a member
// record in the database that matches the provided workspace ID and account ID.
// Each member is uniquely identified by workspace ID and account ID.
func (s *AuthService) AuthenticateWorkspaceMember(
	ctx context.Context, workspaceId string, accountId string,
) (models.Member, error) {
	member, err := s.memberRepo.LookupByWorkspaceAccountId(ctx, workspaceId, accountId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}
	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

func (s *AuthService) ValidatePersonalAccessToken(
	ctx context.Context, token string) (models.Account, error) {
	account, err := s.accountRepo.LookupByToken(ctx, token)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Account{}, ErrAccountNotFound
	}

	if err != nil {
		return models.Account{}, ErrAccount

	}

	return account, nil
}

type CustomerAuthService struct {
	repo ports.CustomerRepositorer
}

func NewCustomerAuthService(repo ports.CustomerRepositorer) *CustomerAuthService {
	return &CustomerAuthService{
		repo: repo,
	}
}

func (s *CustomerAuthService) AuthenticateWorkspaceCustomer(
	ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error) {
	customer, err := s.repo.LookupWorkspaceCustomerById(ctx, workspaceId, customerId, role)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrCustomerNotFound
	}

	if err != nil {
		return models.Customer{}, ErrCustomer
	}
	return customer, nil
}

func (s *CustomerAuthService) GetWidgetLinkedSecretKey(
	ctx context.Context, widgetId string) (models.WorkspaceSecret, error) {
	sk, err := s.repo.LookupSecretKeyByWidgetId(ctx, widgetId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.WorkspaceSecret{}, ErrSecretKeyNotFound
	}

	if err != nil {
		return models.WorkspaceSecret{}, ErrSecretKey
	}
	return sk, nil
}
