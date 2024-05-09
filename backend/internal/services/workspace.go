package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
)

type WorkspaceService struct {
	workspaceRepo ports.WorkspaceRepositorer
	memberRepo    ports.MemberRepositorer
}

func NewWorkspaceService(workspaceRepo ports.WorkspaceRepositorer, memberRepo ports.MemberRepositorer) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		memberRepo:    memberRepo,
	}
}

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, w domain.Workspace) (domain.Workspace, error) {
	workspace, err := s.workspaceRepo.CreateWorkspace(ctx, w)

	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return workspace, ErrWorkspace
	}

	if err != nil {
		return workspace, err
	}
	// TODO: probably send a mail to notify that a new workspace was created
	return workspace, nil
}

func (s *WorkspaceService) GetUserWorkspace(ctx context.Context, accountId string, workspaceId string) (domain.Workspace, error) {
	workspace, err := s.workspaceRepo.GetByAccountWorkspaceId(ctx, accountId, workspaceId)

	if errors.Is(err, repository.ErrQuery) {
		return workspace, ErrWorkspace
	}

	if errors.Is(err, repository.ErrEmpty) {
		return workspace, ErrWorkspaceNotFound
	}

	if err != nil {
		return workspace, err
	}
	return workspace, nil
}

func (s *WorkspaceService) GetUserWorkspaceList(ctx context.Context, accountId string) ([]domain.Workspace, error) {
	workspaces, err := s.workspaceRepo.GetListByAccountId(ctx, accountId)

	if errors.Is(err, repository.ErrQuery) {
		return workspaces, ErrWorkspace
	}

	if err != nil {
		return workspaces, err
	}
	// TODO: add pagination support
	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspaceMember(ctx context.Context, accountId string, workspaceId string) (domain.Member, error) {
	member, err := s.memberRepo.GetByAccountWorkspaceId(ctx, accountId, workspaceId)
	if err != nil {
		return member, err
	}
	return member, nil
}

func (s *WorkspaceService) InitWorkspaceLabel(ctx context.Context, label domain.Label) (domain.Label, bool, error) {
	label, created, err := s.workspaceRepo.GetOrCreateLabel(ctx, label)
	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return domain.Label{}, false, ErrLabel
	}

	if err != nil {
		return domain.Label{}, false, ErrLabel
	}

	return label, created, err
}
