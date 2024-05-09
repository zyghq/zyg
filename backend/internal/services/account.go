package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
)

type AccountService struct {
	repo ports.AccountRepositorer
}

func NewAccountService(repo ports.AccountRepositorer) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) InitiateAccount(ctx context.Context, a domain.Account) (domain.Account, bool, error) {
	account, created, err := s.repo.GetOrCreateByAuthUserId(ctx, a)

	// checks if the result was empty or have query error
	if errors.Is(err, repository.ErrEmpty) || errors.Is(err, repository.ErrQuery) {
		return account, created, ErrAccount
	}

	if err != nil {
		return account, created, err
	}

	return account, created, nil
}

func (s *AccountService) GetAuthUser(ctx context.Context, authUserId string) (domain.Account, error) {
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

func (s *AccountService) IssuePersonalAccessToken(ctx context.Context, ap domain.AccountPAT) (domain.AccountPAT, error) {
	ap, err := s.repo.CreatePersonalAccessToken(ctx, ap)

	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return ap, ErrPat
	}

	if err != nil {
		return ap, err
	}
	// probably send a mail that a new token was created
	// send via background job
	return ap, nil
}

func (s *AccountService) GetUserPatList(ctx context.Context, accountId string) ([]domain.AccountPAT, error) {
	apList, err := s.repo.GetPatListByAccountId(ctx, accountId)

	if errors.Is(err, repository.ErrQuery) {
		return apList, ErrPat
	}

	if err != nil {
		return apList, err
	}

	return apList, nil
}

func (s *AccountService) GetPatAccount(ctx context.Context, token string) (domain.Account, error) {
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
