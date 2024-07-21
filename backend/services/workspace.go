package services

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

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

func (s *WorkspaceService) CreateWorkspace(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.InsertWorkspaceForAccount(ctx, a, w)
	if err != nil {
		return workspace, err
	}
	// TODO: do some engagement here, once the workspace is created.
	return workspace, nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, w models.Workspace) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.ModifyWorkspaceById(ctx, w)
	if err != nil {
		return workspace, err
	}
	return workspace, nil
}

func (s *WorkspaceService) SetWorkspaceLabel(ctx context.Context, workspaceId string, label models.Label) (models.Label, error) {
	label, err := s.workspaceRepo.AlterWorkspaceLabelById(ctx, workspaceId, label)
	if err != nil {
		return label, err
	}

	return label, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.FetchWorkspaceById(ctx, workspaceId)

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

func (s *WorkspaceService) GetMemberWorkspace(
	ctx context.Context, accountId string, workspaceId string,
) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.RetrieveByAccountWorkspaceId(ctx, accountId, workspaceId)

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

func (s *WorkspaceService) ListMemberWorkspaces(ctx context.Context, accountId string) ([]models.Workspace, error) {
	workspaces, err := s.workspaceRepo.FetchWorkspacesByMemberAccountId(ctx, accountId)

	if errors.Is(err, repository.ErrQuery) {
		return workspaces, ErrWorkspace
	}

	if err != nil {
		return workspaces, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) GetWorkspaceMember(ctx context.Context, accountId string, workspaceId string) (models.Member, error) {
	member, err := s.memberRepo.LookupByAccountWorkspaceId(ctx, accountId, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return member, ErrMemberNotFound
	}

	if err != nil {
		return member, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) ListWorkspaceMembers(ctx context.Context, workspaceId string) ([]models.Member, error) {
	members, err := s.memberRepo.RetrieveMembersByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return members, ErrMember
	}
	return members, nil
}

func (s *WorkspaceService) GetWorkspaceMemberById(ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	member, err := s.memberRepo.FetchByWorkspaceMemberId(ctx, workspaceId, memberId)

	if errors.Is(err, repository.ErrEmpty) {
		return member, ErrMemberNotFound
	}

	if err != nil {
		return member, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) ListWorkspaceCustomers(ctx context.Context, workspaceId string) ([]models.Customer, error) {
	customers, err := s.customerRepo.FetchCustomersByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return customers, ErrCustomer
	}
	return customers, nil
}

func (s *WorkspaceService) CreateLabel(ctx context.Context, label models.Label) (models.Label, bool, error) {
	label, created, err := s.workspaceRepo.UpsertLabel(ctx, label)
	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return models.Label{}, false, ErrLabel
	}

	if err != nil {
		return models.Label{}, false, ErrLabel
	}

	return label, created, err
}

func (s *WorkspaceService) GetWorkspaceLabel(ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	label, err := s.workspaceRepo.LookupWorkspaceLabelById(ctx, workspaceId, labelId)

	if errors.Is(err, repository.ErrQuery) {
		return label, ErrLabel
	}

	if errors.Is(err, repository.ErrEmpty) {
		return label, ErrLabelNotFound
	}
	return label, err
}

func (s *WorkspaceService) ListWorkspaceLabels(ctx context.Context, workspaceId string) ([]models.Label, error) {
	labels, err := s.workspaceRepo.RetrieveLabelsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return labels, ErrLabel
	}
	return labels, nil
}

func (s *WorkspaceService) CreateCustomerWithExternalId(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	customer, created, err := s.customerRepo.UpsertCustomerByExtId(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateCustomerWithEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	customer, created, err := s.customerRepo.UpsertCustomerByEmail(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateCustomerWithPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	customer, created, err := s.customerRepo.UpsertCustomerByPhone(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateAnonymousCustomer(ctx context.Context, c models.Customer) (models.Customer, bool, error) {
	customer, created, err := s.customerRepo.UpsertCustomerByAnonId(ctx, c)
	if err != nil {
		return customer, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) AddMember(ctx context.Context, workspace models.Workspace, member models.Member) (models.Member, error) {
	member, err := s.workspaceRepo.InsertMemberIntoWorkspace(ctx, workspace.WorkspaceId, member)
	if err != nil {
		return member, err
	}
	return member, nil
}

func (s *WorkspaceService) CreateWidget(ctx context.Context, workspaceId string, widget models.Widget) (models.Widget, error) {
	widget, err := s.workspaceRepo.InsertWidgetIntoWorkspace(ctx, workspaceId, widget)
	if err != nil {
		return widget, err
	}
	return widget, nil
}

func (s *WorkspaceService) ListWidgets(ctx context.Context, workspaceId string) ([]models.Widget, error) {
	widgets, err := s.workspaceRepo.RetrieveWidgetsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return widgets, err
	}
	return widgets, nil
}

func (s *WorkspaceService) GenerateSecretKey(ctx context.Context, workspaceId string, length int) (models.SecretKey, error) {
	// Create a buffer to store our entropy sources
	entropy := make([]byte, 0, 1024)
	// Add current timestamp
	entropy = append(entropy, []byte(time.Now().String())...)
	// Add process ID
	entropy = append(entropy, []byte(fmt.Sprintf("%d", os.Getpid()))...)
	// Add random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return models.SecretKey{}, err
	}

	entropy = append(entropy, randomBytes...)

	// Hash the entropy using SHA-512
	hash := sha512.Sum512(entropy)
	// Encode the hash using base64
	encodedHash := base64.URLEncoding.EncodeToString(hash[:])
	// Return the first 'length' characters of the encoded hash
	slicedHash := encodedHash[:length]

	sk, err := s.workspaceRepo.InsertSecretKeyIntoWorkspace(ctx, workspaceId, slicedHash)
	if err != nil {
		return sk, err
	}
	return sk, nil
}

func (s *WorkspaceService) GetWorkspaceSecretKey(ctx context.Context, workspaceId string) (models.SecretKey, error) {
	sk, err := s.workspaceRepo.FetchSecretKeyByWorkspaceId(ctx, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return sk, ErrSecretKeyNotFound
	}

	if err != nil {
		return sk, ErrSecretKey
	}
	return sk, nil
}

func (s *WorkspaceService) GetWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error) {
	widget, err := s.workspaceRepo.LookupWorkspaceWidget(ctx, widgetId)

	if errors.Is(err, repository.ErrQuery) {
		return widget, ErrWidget
	}

	if errors.Is(err, repository.ErrEmpty) {
		return widget, ErrWidgetNotFound
	}

	if err != nil {
		return widget, err
	}
	return widget, nil
}
