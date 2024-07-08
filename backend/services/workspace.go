package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
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
		customerRepo:  customerRepo,
	}
}

func (s *WorkspaceService) CreateAccountWorkspace(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.CreateWorkspaceByAccount(ctx, a, w)

	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return workspace, ErrWorkspace
	}

	if err != nil {
		return workspace, err
	}
	// TODO: probably send a mail to notify that a new workspace was created
	return workspace, nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, w models.Workspace) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.UpdateWorkspaceById(ctx, w)
	if err != nil {
		return workspace, err
	}
	return workspace, nil
}

func (s *WorkspaceService) UpdateWorkspaceLabel(ctx context.Context, workspaceId string, label models.Label) (models.Label, error) {
	label, err := s.workspaceRepo.UpdateWorkspaceLabelById(ctx, workspaceId, label)
	if err != nil {
		return label, err
	}

	return label, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error) {
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

func (s *WorkspaceService) MemberWorkspace(
	ctx context.Context, accountId string, workspaceId string,
) (models.Workspace, error) {
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

func (s *WorkspaceService) MemberWorkspaces(ctx context.Context, accountId string) ([]models.Workspace, error) {
	workspaces, err := s.workspaceRepo.GetListByMemberAccountId(ctx, accountId)

	if errors.Is(err, repository.ErrQuery) {
		return workspaces, ErrWorkspace
	}

	if err != nil {
		return workspaces, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) WorkspaceUserMember(ctx context.Context, accountId string, workspaceId string) (models.Member, error) {
	member, err := s.memberRepo.GetByAccountWorkspaceId(ctx, accountId, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return member, ErrMemberNotFound
	}

	if err != nil {
		return member, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) WorkspaceMembers(ctx context.Context, workspaceId string) ([]models.Member, error) {
	members, err := s.memberRepo.GetListByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return members, ErrMember
	}
	return members, nil
}

func (s *WorkspaceService) WorkspaceMember(ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	member, err := s.memberRepo.GetByWorkspaceMemberId(ctx, workspaceId, memberId)

	if errors.Is(err, repository.ErrEmpty) {
		return member, ErrMemberNotFound
	}

	if err != nil {
		return member, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) WorkspaceCustomers(ctx context.Context, workspaceId string) ([]models.Customer, error) {
	customers, err := s.customerRepo.GetListByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return customers, ErrCustomer
	}
	return customers, nil
}

func (s *WorkspaceService) InitWorkspaceLabel(ctx context.Context, label models.Label) (models.Label, bool, error) {
	label, created, err := s.workspaceRepo.GetOrCreateLabel(ctx, label)
	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return models.Label{}, false, ErrLabel
	}

	if err != nil {
		return models.Label{}, false, ErrLabel
	}

	return label, created, err
}

func (s *WorkspaceService) WorkspaceLabel(ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	label, err := s.workspaceRepo.GetWorkspaceLabelById(ctx, workspaceId, labelId)

	if errors.Is(err, repository.ErrQuery) {
		return label, ErrLabel
	}

	if errors.Is(err, repository.ErrEmpty) {
		return label, ErrLabelNotFound
	}
	return label, err
}

func (s *WorkspaceService) WorkspaceLabels(ctx context.Context, workspaceId string) ([]models.Label, error) {
	labels, err := s.workspaceRepo.GetLabelListByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return labels, ErrLabel
	}
	return labels, nil
}

func (s *WorkspaceService) InitWorkspaceCustomerWithExternalId(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	defaultName := c.AnonName()
	if !c.Name.Valid {
		c.Name = models.NullString(&defaultName)
	}

	customer, created, err := s.customerRepo.GetOrCreateCustomerByExtId(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) InitWorkspaceCustomerWithEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	defaultName := c.AnonName()
	if !c.Name.Valid {
		c.Name = models.NullString(&defaultName)
	}

	customer, created, err := s.customerRepo.GetOrCreateCustomerByEmail(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) InitWorkspaceCustomerWithPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	defaultName := c.AnonName()
	if !c.Name.Valid {
		c.Name = models.NullString(&defaultName)
	}

	customer, created, err := s.customerRepo.GetOrCreateCustomerByPhone(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) AddMember(ctx context.Context, workspace models.Workspace, member models.Member) (models.Member, error) {
	member, err := s.workspaceRepo.AddMemberByWorkspaceId(ctx, workspace.WorkspaceId, member)
	if err != nil {
		return member, err
	}
	return member, nil
}
