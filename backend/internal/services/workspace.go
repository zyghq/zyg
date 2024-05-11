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
	customerRepo  ports.CustomerRepositorer
}

func NewWorkspaceService(
	workspaceRepo ports.WorkspaceRepositorer,
	memberRepo ports.MemberRepositorer,
	customerRepo ports.CustomerRepositorer,
) *WorkspaceService {
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

func (s *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId string) (domain.Workspace, error) {
	workspace, err := s.workspaceRepo.GetWorkspaceById(ctx, workspaceId)

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

func (s *WorkspaceService) GetUserWorkspace(
	ctx context.Context, accountId string, workspaceId string,
) (domain.Workspace, error) {
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

	if errors.Is(err, repository.ErrEmpty) {
		return member, ErrMemberNotFound
	}

	if err != nil {
		return member, ErrMember
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

func (s *WorkspaceService) InitWorkspaceCustomerWithExternalId(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.customerRepo.GetOrCreateCustomerByExtId(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) InitWorkspaceCustomerWithEmail(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.customerRepo.GetOrCreateCustomerByEmail(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) InitWorkspaceCustomerWithPhone(ctx context.Context, c domain.Customer) (domain.Customer, bool, error) {
	customer, created, err := s.customerRepo.GetOrCreateCustomerByPhone(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}
