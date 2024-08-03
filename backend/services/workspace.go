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

func (s *WorkspaceService) CreateWorkspace(
	ctx context.Context, accountId string, memberName string, workspaceName string) (models.Workspace, error) {
	workspace := models.Workspace{
		AccountId: accountId,
		Name:      workspaceName,
	}
	workspace, err := s.workspaceRepo.InsertWorkspace(ctx, memberName, workspace)
	if err != nil {
		return models.Workspace{}, err
	}
	return workspace, nil
}

func (s *WorkspaceService) UpdateWorkspace(
	ctx context.Context, workspace models.Workspace) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.ModifyWorkspaceById(ctx, workspace)
	if err != nil {
		return models.Workspace{}, err
	}
	return workspace, nil
}

func (s *WorkspaceService) UpdateLabel(
	ctx context.Context, label models.Label) (models.Label, error) {
	label, err := s.workspaceRepo.ModifyLabelById(ctx, label)
	if err != nil {
		return models.Label{}, err
	}

	return label, nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error) {
	workspace, err := s.workspaceRepo.FetchByWorkspaceId(ctx, workspaceId)

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

func (s *WorkspaceService) GetAccountLinkedWorkspace(
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

func (s *WorkspaceService) ListAccountLinkedWorkspaces(
	ctx context.Context, accountId string) ([]models.Workspace, error) {
	workspaces, err := s.workspaceRepo.FetchWorkspacesByAccountId(ctx, accountId)
	if err != nil {
		return []models.Workspace{}, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) GetAccountLinkedMember(
	ctx context.Context, workspaceId string, accountId string) (models.Member, error) {
	member, err := s.memberRepo.LookupByWorkspaceAccountId(ctx, workspaceId, accountId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}

	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) ListMembers(
	ctx context.Context, workspaceId string) ([]models.Member, error) {
	members, err := s.memberRepo.RetrieveMembersByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Member{}, ErrMember
	}
	return members, nil
}

func (s *WorkspaceService) GetMember(
	ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	member, err := s.memberRepo.FetchByWorkspaceMemberId(ctx, workspaceId, memberId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}

	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

func (s *WorkspaceService) ListCustomers(
	ctx context.Context, workspaceId string) ([]models.Customer, error) {
	customers, err := s.customerRepo.FetchCustomersByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Customer{}, err
	}
	return customers, nil
}

func (s *WorkspaceService) CreateLabel(
	ctx context.Context, workspaceId string, name string, icon string) (models.Label, bool, error) {
	label := models.Label{
		WorkspaceId: workspaceId,
		Name:        name,
		Icon:        icon,
	}
	label, created, err := s.workspaceRepo.InsertLabelByName(ctx, label)
	if err != nil {
		return models.Label{}, false, ErrLabel
	}

	return label, created, err
}

func (s *WorkspaceService) GetLabel(
	ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	label, err := s.workspaceRepo.LookupWorkspaceLabelById(ctx, workspaceId, labelId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Label{}, ErrLabelNotFound
	}

	if err != nil {
		return models.Label{}, ErrLabel
	}

	return label, err
}

func (s *WorkspaceService) ListLabels(
	ctx context.Context, workspaceId string) ([]models.Label, error) {
	labels, err := s.workspaceRepo.FetchLabelsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Label{}, ErrLabel
	}
	return labels, nil
}

func (s *WorkspaceService) CreateCustomerWithExternalId(
	ctx context.Context, workspaceId string, externalId string, isVerified bool, name string,
) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		ExternalId:  models.NullString(&externalId),
		IsVerified:  isVerified,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := s.customerRepo.UpsertCustomerByExtId(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateCustomerWithEmail(
	ctx context.Context, workspaceId string, email string, isVerified bool, name string) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		Email:       models.NullString(&email),
		IsVerified:  isVerified,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := s.customerRepo.UpsertCustomerByEmail(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateCustomerWithPhone(
	ctx context.Context, workspaceId string, phone string, isVerified bool, name string) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		Phone:       models.NullString(&phone),
		IsVerified:  isVerified,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := s.customerRepo.UpsertCustomerByPhone(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (s *WorkspaceService) CreateAnonymousCustomer(
	ctx context.Context, workspaceId string, anonId string, isVerified bool, name string) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		AnonId:      anonId,
		IsVerified:  isVerified,
		Name:        name,
		Role:        models.Customer{}.Visitor(),
	}
	customer, created, err := s.customerRepo.UpsertCustomerByAnonId(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

// func (s *WorkspaceService) AddMember(
// 	ctx context.Context, workspaceId string, member models.Member) (models.Member, error) {
// 	member, err := s.workspaceRepo.InsertMember(ctx, workspaceId, member)
// 	if err != nil {
// 		return member, err
// 	}
// 	return member, nil
// }

func (s *WorkspaceService) CreateWidget(
	ctx context.Context, workspaceId string, name string, configuration map[string]interface{}) (models.Widget, error) {
	widget := models.Widget{
		WorkspaceId:   workspaceId,
		Name:          name,
		Configuration: configuration,
	}
	widget, err := s.workspaceRepo.InsertWidget(ctx, widget)
	if err != nil {
		return models.Widget{}, err
	}
	return widget, nil
}

func (s *WorkspaceService) ListWidgets(
	ctx context.Context, workspaceId string) ([]models.Widget, error) {
	widgets, err := s.workspaceRepo.FetchWidgetsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Widget{}, err
	}
	return widgets, nil
}

func (s *WorkspaceService) GenerateSecretKey(
	ctx context.Context, workspaceId string, length int) (models.SecretKey, error) {
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
	slicedHash := "sk" + encodedHash[:length]

	sk, err := s.workspaceRepo.InsertSecretKey(ctx, workspaceId, slicedHash)
	if err != nil {
		return models.SecretKey{}, err
	}
	return sk, nil
}

func (s *WorkspaceService) GetSecretKey(
	ctx context.Context, workspaceId string) (models.SecretKey, error) {
	sk, err := s.workspaceRepo.FetchSecretKeyByWorkspaceId(ctx, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.SecretKey{}, ErrSecretKeyNotFound
	}

	if err != nil {
		return models.SecretKey{}, ErrSecretKey
	}
	return sk, nil
}

func (s *WorkspaceService) GetWidget(
	ctx context.Context, widgetId string) (models.Widget, error) {
	widget, err := s.workspaceRepo.LookupWidgetById(ctx, widgetId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Widget{}, ErrWidgetNotFound
	}

	if err != nil {
		return models.Widget{}, err
	}

	return widget, nil
}
