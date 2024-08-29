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

func (ws *WorkspaceService) UpdateWorkspace(
	ctx context.Context, workspace models.Workspace) (models.Workspace, error) {
	workspace, err := ws.workspaceRepo.ModifyWorkspaceById(ctx, workspace)
	if err != nil {
		return models.Workspace{}, err
	}
	return workspace, nil
}

func (ws *WorkspaceService) UpdateLabel(
	ctx context.Context, label models.Label) (models.Label, error) {
	label, err := ws.workspaceRepo.ModifyLabelById(ctx, label)
	if err != nil {
		return models.Label{}, err
	}

	return label, nil
}

func (ws *WorkspaceService) GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error) {
	workspace, err := ws.workspaceRepo.FetchByWorkspaceId(ctx, workspaceId)

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

func (ws *WorkspaceService) GetAccountLinkedMember(
	ctx context.Context, workspaceId string, accountId string) (models.Member, error) {
	member, err := ws.memberRepo.LookupByWorkspaceAccountId(ctx, workspaceId, accountId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}

	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

func (ws *WorkspaceService) ListMembers(
	ctx context.Context, workspaceId string) ([]models.Member, error) {
	members, err := ws.memberRepo.FetchMembersByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Member{}, ErrMember
	}
	return members, nil
}

func (ws *WorkspaceService) GetMember(
	ctx context.Context, workspaceId string, memberId string) (models.Member, error) {
	member, err := ws.memberRepo.FetchByWorkspaceMemberId(ctx, workspaceId, memberId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}

	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

func (ws *WorkspaceService) ListCustomers(
	ctx context.Context, workspaceId string) ([]models.Customer, error) {
	role := models.Customer{}.Engaged()
	customers, err := ws.customerRepo.FetchCustomersByWorkspaceId(ctx, workspaceId, &role)
	if err != nil {
		return []models.Customer{}, err
	}
	return customers, nil
}

func (ws *WorkspaceService) CreateLabel(
	ctx context.Context, workspaceId string, name string, icon string) (models.Label, bool, error) {
	label := models.Label{
		WorkspaceId: workspaceId,
		Name:        name,
		Icon:        icon,
	}
	label, created, err := ws.workspaceRepo.InsertLabelByName(ctx, label)
	if err != nil {
		return models.Label{}, false, ErrLabel
	}

	return label, created, err
}

func (ws *WorkspaceService) GetLabel(
	ctx context.Context, workspaceId string, labelId string) (models.Label, error) {
	label, err := ws.workspaceRepo.LookupWorkspaceLabelById(ctx, workspaceId, labelId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Label{}, ErrLabelNotFound
	}

	if err != nil {
		return models.Label{}, ErrLabel
	}

	return label, err
}

func (ws *WorkspaceService) ListLabels(
	ctx context.Context, workspaceId string) ([]models.Label, error) {
	labels, err := ws.workspaceRepo.FetchLabelsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Label{}, ErrLabel
	}
	return labels, nil
}

func (ws *WorkspaceService) CreateCustomerWithExternalId(
	ctx context.Context, workspaceId string, externalId string, isAnonymous bool, name string,
) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		ExternalId:  models.NullString(&externalId),
		IsAnonymous: isAnonymous,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByExtId(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateCustomerWithEmail(
	ctx context.Context, workspaceId string, email string, isAnonymous bool, name string) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		Email:       models.NullString(&email),
		IsAnonymous: isAnonymous,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByEmail(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateCustomerWithPhone(
	ctx context.Context, workspaceId string, phone string, isAnonymous bool, name string) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customer := models.Customer{
		WorkspaceId: workspaceId,
		Phone:       models.NullString(&phone),
		IsAnonymous: isAnonymous,
		Name:        name,
		Role:        models.Customer{}.Engaged(),
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByPhone(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateAnonymousCustomer(
	ctx context.Context, workspaceId string, name string) (models.Customer, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customerId := models.Customer{}.GenId()
	customer := models.Customer{
		CustomerId:  customerId,
		WorkspaceId: workspaceId,
		IsAnonymous: true,
		Name:        name,
		Role:        models.Customer{}.Visitor(),
	}
	// Modify, insert function, just save it to the db.
	customer, _, err := ws.customerRepo.UpsertCustomerById(ctx, customer)
	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (ws *WorkspaceService) CreateWidget(
	ctx context.Context, workspaceId string, name string,
	configuration map[string]interface{}) (models.Widget, error) {
	widget := models.Widget{
		WorkspaceId:   workspaceId,
		Name:          name,
		Configuration: configuration,
	}
	widget, err := ws.workspaceRepo.InsertWidget(ctx, widget)
	if err != nil {
		return models.Widget{}, err
	}
	return widget, nil
}

func (ws *WorkspaceService) ListWidgets(
	ctx context.Context, workspaceId string) ([]models.Widget, error) {
	widgets, err := ws.workspaceRepo.FetchWidgetsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return []models.Widget{}, err
	}
	return widgets, nil
}

func (ws *WorkspaceService) GenerateWorkspaceSecret(
	ctx context.Context, workspaceId string, length int) (models.WorkspaceSecret, error) {
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
		return models.WorkspaceSecret{}, err
	}

	entropy = append(entropy, randomBytes...)

	// Hash the entropy using SHA-512
	hash := sha512.Sum512(entropy)
	// Encode the hash using base64
	encodedHash := base64.URLEncoding.EncodeToString(hash[:])
	// Return the first 'length' characters of the encoded hash
	slicedHash := "sk" + encodedHash[:length]

	sk, err := ws.workspaceRepo.InsertWorkspaceSecret(ctx, workspaceId, slicedHash)
	if err != nil {
		return models.WorkspaceSecret{}, err
	}
	return sk, nil
}

func (ws *WorkspaceService) GetSecretKey(
	ctx context.Context, workspaceId string) (models.WorkspaceSecret, error) {
	sk, err := ws.workspaceRepo.FetchSecretKeyByWorkspaceId(ctx, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.WorkspaceSecret{}, ErrSecretKeyNotFound
	}

	if err != nil {
		return models.WorkspaceSecret{}, ErrSecretKey
	}
	return sk, nil
}

func (ws *WorkspaceService) GetWidget(
	ctx context.Context, widgetId string) (models.Widget, error) {
	widget, err := ws.workspaceRepo.LookupWidgetById(ctx, widgetId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Widget{}, ErrWidgetNotFound
	}
	if err != nil {
		return models.Widget{}, ErrWidget
	}
	return widget, nil
}

func (ws *WorkspaceService) GetCustomer(
	ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error) {
	customer, err := ws.customerRepo.LookupWorkspaceCustomerById(ctx, workspaceId, customerId, role)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrCustomerNotFound
	}
	if err != nil {
		return models.Customer{}, ErrCustomer
	}
	return customer, nil
}

func (ws *WorkspaceService) DoesEmailConflict(
	ctx context.Context, workspaceId string, email string) (bool, error) {
	exists, err := ws.customerRepo.CheckEmailExists(ctx, workspaceId, email)
	if err != nil {
		// Be pessimistic when checking for email conflict, let us assume email already exists;
		// hence there is conflict.
		return true, ErrCustomer
	}
	return exists, nil
}

func (ws *WorkspaceService) ValidateWidgetSession(
	ctx context.Context, sk string, widgetId string, sessionId string) (models.Customer, error) {
	// fetch widget session for the provided widget session ID.
	session, err := ws.workspaceRepo.LookupWidgetSessionById(ctx, widgetId, sessionId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrWidgetSessionInvalid
	}
	if err != nil {
		return models.Customer{}, ErrWidgetSession
	}

	// decode data from the widget session.
	// if there is an error, we assume the widget session is invalid.
	data, err := session.Decode(sk)
	if err != nil {
		fmt.Println(err)
		return models.Customer{}, ErrWidgetSessionInvalid
	}

	customer, err := ws.customerRepo.LookupWorkspaceCustomerById(ctx, data.WorkspaceId, data.CustomerId, nil)
	if err != nil {
		return models.Customer{}, ErrWidgetSessionInvalid
	}
	// check the calculated identity hash against the one in the widget session.
	if customer.IdentityHash() != data.IdentityHash {
		return models.Customer{}, ErrWidgetSessionInvalid
	}
	return customer, nil
}

func (ws *WorkspaceService) CreateWidgetSession(
	ctx context.Context, sk string, workspaceId string, widgetId string,
	sessionId string, name string) (models.Customer, bool, error) {
	var created bool

	// create a new customer with the provided name.
	customer, err := ws.CreateAnonymousCustomer(ctx, workspaceId, name)
	if err != nil {
		return models.Customer{}, created, ErrWidgetSession
	}

	// creates a new widget session with session data for the new customer ID
	// and calculated identity hash.
	session := (&models.WidgetSession{}).CreateSession(sessionId, widgetId)
	data := session.CreateSessionData(workspaceId, customer.CustomerId, customer.IdentityHash())
	// set the encoded data for the provided secret key.
	err = session.SetEncodeData(sk, data)
	if err != nil {
		return models.Customer{}, created, ErrWidgetSession
	}
	// insert the new widget session into the db.
	_, created, err = ws.workspaceRepo.UpsertWidgetSessionById(ctx, session)
	if err != nil {
		return models.Customer{}, created, ErrWidgetSession
	}
	return customer, created, nil
}
