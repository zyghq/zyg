package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type AccountService struct {
	accountRepo   ports.AccountRepositorer
	workspaceRepo ports.WorkspaceRepositorer
}

func NewAccountService(
	accountRepo ports.AccountRepositorer,
	workspaceRepo ports.WorkspaceRepositorer) *AccountService {
	return &AccountService{
		accountRepo:   accountRepo,
		workspaceRepo: workspaceRepo,
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
	account, created, err := s.accountRepo.UpsertByAuthUserId(ctx, account)
	if err != nil {
		return models.Account{}, false, ErrAccount
	}
	return account, created, nil
}

func (s *AccountService) GeneratePersonalAccessToken(
	ctx context.Context, accountId string, name string, description string) (models.AccountPAT, error) {
	pat := models.AccountPAT{
		AccountId:   accountId,
		Name:        name,
		UnMask:      true, // unmask only when created
		Description: description,
	}
	ap, err := s.accountRepo.InsertPersonalAccessToken(ctx, pat)
	if err != nil {
		return models.AccountPAT{}, ErrPat
	}
	// @sanchitrk
	// send an email that a new token was created.
	return ap, nil
}

func (s *AccountService) ListPersonalAccessTokens(
	ctx context.Context, accountId string) ([]models.AccountPAT, error) {
	pats, err := s.accountRepo.FetchPatsByAccountId(ctx, accountId)
	if err != nil {
		return []models.AccountPAT{}, ErrPat
	}

	return pats, nil
}

func (s *AccountService) GetPersonalAccessToken(
	ctx context.Context, patId string) (models.AccountPAT, error) {
	pat, err := s.accountRepo.FetchPatById(ctx, patId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.AccountPAT{}, ErrPatNotFound
	}

	if err != nil {
		return models.AccountPAT{}, ErrPat
	}

	return pat, nil
}

func (s *AccountService) DeletePersonalAccessToken(ctx context.Context, patId string) error {
	err := s.accountRepo.DeletePatById(ctx, patId)
	if err != nil {
		return ErrPat
	}
	// @sanchitrk
	// send an email that the token is deleted
	return nil
}

func (s *AccountService) CreateWorkspace(
	ctx context.Context, accountId string, memberName string, workspaceName string) (models.Workspace, error) {
	workspace := models.Workspace{
		AccountId: accountId,
		Name:      workspaceName,
	}
	member := models.Member{
		Name: memberName,
		Role: models.MemberRole{}.Primary(),
	}
	workspace, err := s.workspaceRepo.InsertWorkspaceWithMember(ctx, workspace, member)
	if err != nil {
		return models.Workspace{}, err
	}
	return workspace, nil
}

func (s *AccountService) GetAccountLinkedWorkspace(
	ctx context.Context, accountId string, workspaceId string) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.LookupWorkspaceByAccountId(ctx, workspaceId, accountId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Workspace{}, ErrWorkspaceNotFound
	}

	if err != nil {
		return models.Workspace{}, err
	}
	return workspace, nil
}

func (s *AccountService) ListAccountLinkedWorkspaces(
	ctx context.Context, accountId string) ([]models.Workspace, error) {
	workspaces, err := s.workspaceRepo.FetchWorkspacesByAccountId(ctx, accountId)
	if err != nil {
		return []models.Workspace{}, err
	}
	return workspaces, nil
}
