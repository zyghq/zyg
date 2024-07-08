package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type AccountService struct {
	repo ports.AccountRepositorer
}

func NewAccountService(repo ports.AccountRepositorer) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) InitiateAccount(ctx context.Context, a models.Account) (models.Account, bool, error) {
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

func (s *AccountService) AuthenticateUser(ctx context.Context, authUserId string) (models.Account, error) {
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

func (s *AccountService) GeneratePersonalAccessToken(ctx context.Context, ap models.AccountPAT) (models.AccountPAT, error) {
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

func (s *AccountService) GetPersonalAccessTokens(ctx context.Context, accountId string) ([]models.AccountPAT, error) {
	pats, err := s.repo.GetPatListByAccountId(ctx, accountId)

	if errors.Is(err, repository.ErrQuery) {
		return pats, ErrPat
	}

	if err != nil {
		return pats, err
	}

	return pats, nil
}

func (s *AccountService) GetPersonalAccessToken(ctx context.Context, patId string) (models.AccountPAT, error) {
	pat, err := s.repo.GetPatByPatId(ctx, patId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.AccountPAT{}, ErrPatNotFound
	}

	if errors.Is(err, repository.ErrQuery) {
		return models.AccountPAT{}, ErrPat
	}

	if err != nil {
		return models.AccountPAT{}, err
	}

	return pat, nil
}

func (s *AccountService) DeletePersonalAccessToken(ctx context.Context, patId string) error {
	err := s.repo.HardDeletePatByPatId(ctx, patId)
	if err != nil {
		return err
	}
	// probably send a mail that the token was deleted
	// send via background job
	return nil
}
