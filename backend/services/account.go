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

func (s *AccountService) CreateAuthAccount(
	ctx context.Context, authUserId string, email string, name string, provider string) (models.Account, bool, error) {
	account := models.Account{
		AuthUserId: authUserId,
		Email:      email,
		Name:       name,
		Provider:   provider,
	}
	account, created, err := s.repo.UpsertByAuthUserId(ctx, account)
	if err != nil {
		return models.Account{}, false, ErrAccount
	}
	return account, created, nil
}

func (s *AccountService) AuthenticateUser(
	ctx context.Context, authUserId string) (models.Account, error) {
	account, err := s.repo.FetchByAuthUserId(ctx, authUserId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Account{}, ErrAccountNotFound
	}

	if err != nil {
		return models.Account{}, ErrAccount
	}

	return account, nil
}

func (s *AccountService) GeneratePersonalAccessToken(
	ctx context.Context, accountId string, name string, description string) (models.AccountPAT, error) {
	pat := models.AccountPAT{
		AccountId:   accountId,
		Name:        name,
		UnMask:      true, // unmask only when created
		Description: description,
	}
	ap, err := s.repo.InsertPersonalAccessToken(ctx, pat)
	if err != nil {
		return models.AccountPAT{}, ErrPat
	}
	// @sanchitrk
	// send an email that a new token was created.
	return ap, nil
}

func (s *AccountService) ListPersonalAccessTokens(
	ctx context.Context, accountId string) ([]models.AccountPAT, error) {
	pats, err := s.repo.FetchPatsByAccountId(ctx, accountId)
	if err != nil {
		return []models.AccountPAT{}, ErrPat
	}

	return pats, nil
}

func (s *AccountService) GetPersonalAccessToken(
	ctx context.Context, patId string) (models.AccountPAT, error) {
	pat, err := s.repo.FetchPatById(ctx, patId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.AccountPAT{}, ErrPatNotFound
	}

	if err != nil {
		return models.AccountPAT{}, ErrPat
	}

	return pat, nil
}

func (s *AccountService) DeletePersonalAccessToken(ctx context.Context, patId string) error {
	err := s.repo.DeletePatById(ctx, patId)
	if err != nil {
		return ErrPat
	}
	// @sanchitrk
	// send an email that the token is deleted
	return nil
}
