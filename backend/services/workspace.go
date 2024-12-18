package services

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/zyghq/postmark"
	"log/slog"
	"os"
	"time"

	"github.com/zyghq/zyg"
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
	workspaceRepo ports.WorkspaceRepositorer, memberRepo ports.MemberRepositorer, customerRepo ports.CustomerRepositorer,
) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepo: workspaceRepo,
		memberRepo:    memberRepo,
		customerRepo:  customerRepo,
	}
}

func (ws *WorkspaceService) GetSystemMember(
	ctx context.Context, workspaceId string) (models.Member, error) {
	member, err := ws.workspaceRepo.LookupSystemMemberByOldest(ctx, workspaceId)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Member{}, ErrMemberNotFound
	}
	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
}

// CreateNewSystemMember creates a new system member for the workspace.
func (ws *WorkspaceService) CreateNewSystemMember(
	ctx context.Context, workspaceId string) (models.Member, error) {
	member := models.Member{}.CreateNewSystemMember(workspaceId)
	member, err := ws.workspaceRepo.InsertSystemMember(ctx, member)
	if err != nil {
		return models.Member{}, ErrMember
	}
	return member, nil
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

func (ws *WorkspaceService) GetWorkspace(
	ctx context.Context, workspaceId string) (models.Workspace, error) {
	hub := sentry.GetHubFromContext(ctx)
	workspace, err := ws.workspaceRepo.FetchByWorkspaceId(ctx, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return workspace, ErrWorkspaceNotFound
	}
	if errors.Is(err, repository.ErrQuery) {
		hub.CaptureException(err)
		return workspace, ErrWorkspace
	}
	if err != nil {
		hub.CaptureException(err)
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
	customers, err := ws.customerRepo.FetchCustomersByWorkspaceId(ctx, workspaceId, nil)
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
	ctx context.Context, workspaceId string, externalId string, name string,
) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	now := time.Now().UTC()
	customer := models.Customer{
		WorkspaceId:     workspaceId,
		ExternalId:      models.NullString(&externalId),
		IsEmailVerified: false, // mark email as unverified.
		Name:            name,
		Role:            models.Customer{}.Engaged(),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByExtId(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateCustomerWithEmail(
	ctx context.Context, workspaceId string, email string, isEmailVerified bool, name string,
) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	now := time.Now().UTC()
	customer := models.Customer{
		WorkspaceId:     workspaceId,
		Email:           models.NullString(&email),
		IsEmailVerified: isEmailVerified,
		Name:            name,
		Role:            models.Customer{}.Engaged(),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByEmail(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateCustomerWithPhone(
	ctx context.Context, workspaceId string, phone string, name string,
) (models.Customer, bool, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	now := time.Now().UTC()
	customer := models.Customer{
		WorkspaceId:     workspaceId,
		Phone:           models.NullString(&phone),
		IsEmailVerified: false, // mark email as unverified
		Name:            name,
		Role:            models.Customer{}.Engaged(),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	customer, created, err := ws.customerRepo.UpsertCustomerByPhone(ctx, customer)
	if err != nil {
		return models.Customer{}, created, err
	}
	return customer, created, nil
}

func (ws *WorkspaceService) CreateUnverifiedCustomer(
	ctx context.Context, workspaceId string, name string) (models.Customer, error) {
	if name == "" {
		name = models.Customer{}.AnonName()
	}
	customerId := models.Customer{}.GenId()
	customer := models.Customer{
		CustomerId:      customerId,
		WorkspaceId:     workspaceId,
		IsEmailVerified: false,
		Name:            name,
		Role:            models.Customer{}.Visitor(),
	}
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
	entropy = append(entropy, []byte(time.Now().UTC().String())...)
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

func (ws *WorkspaceService) GetOrGenerateSecretKey(
	ctx context.Context, workspaceId string) (models.WorkspaceSecret, error) {
	sk, err := ws.GetSecretKey(ctx, workspaceId)
	if errors.Is(err, ErrSecretKeyNotFound) {
		return ws.GenerateWorkspaceSecret(ctx, workspaceId, zyg.DefaultSecretKeyLength)
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
		// any error assume that the email is existing.
		// be pessimistic.
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
		slog.Error("failed to decode session", slog.Any("err", err))
		return models.Customer{}, ErrWidgetSessionInvalid
	}

	customer, err := ws.customerRepo.LookupWorkspaceCustomerById(
		ctx, data.WorkspaceId, data.CustomerId, nil)
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

	// create a new unverified customer with the provided name.
	customer, err := ws.CreateUnverifiedCustomer(ctx, workspaceId, name)
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

func (ws *WorkspaceService) GetCustomerByEmail(
	ctx context.Context, workspaceId string, email string) (models.Customer, error) {
	customer, err := ws.customerRepo.LookupWorkspaceCustomerByEmail(ctx, workspaceId, email, nil)
	if errors.Is(err, repository.ErrEmpty) {
		return models.Customer{}, ErrCustomerNotFound
	}
	if err != nil {
		return models.Customer{}, ErrCustomer
	}
	return customer, nil
}

func (ws *WorkspaceService) PostmarkCreateMailServer(
	ctx context.Context, workspaceId, email, domain string) (models.PostmarkMailServerSetting, error) {
	hub := sentry.GetHubFromContext(ctx)
	// Requires Account Token to manage servers.
	client := postmark.NewClient("", zyg.PostmarkAccountToken())
	// Deprecated:
	// BounceHookUrl
	// OpenHookUrl
	// DeliveryHookUrl
	// ClickHookUrl
	inboundHookURL := fmt.Sprintf("%s/webhooks/%s/postmark/inbound/", zyg.ServerUrl(), workspaceId)
	server := postmark.Server{
		Name:           workspaceId,
		Color:          "Green",
		InboundHookURL: inboundHookURL,
	}
	server, err := client.CreateServer(ctx, server)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}
	// Capture message in Sentry as this helps in keeping track of interactions with Postmark in all env.
	// We also want to make sure to audit.
	hub.CaptureMessage(fmt.Sprintf("postmark server created with ID: %d", server.ID))
	var serverToken string
	if len(server.APITokens) > 0 {
		serverToken = server.APITokens[0]
	}
	now := time.Now().UTC()
	// Save in workspace Postmark setting
	setting := models.PostmarkMailServerSetting{
		WorkspaceId:          workspaceId,
		ServerId:             server.ID,
		ServerToken:          serverToken,
		CreatedAt:            now,
		UpdatedAt:            now,
		IsEnabled:            false,
		Email:                email,
		Domain:               domain,
		HasError:             false,
		InboundEmail:         &server.InboundAddress,
		HasForwardingEnabled: false,
		HasDNS:               false,
		IsDNSVerified:        false,
	}
	setting, err = ws.workspaceRepo.SavePostmarkMailServerSetting(ctx, setting)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}
	return setting, nil
}

func (ws *WorkspaceService) GetPostmarkMailServerSetting(
	ctx context.Context, workspaceId string) (models.PostmarkMailServerSetting, error) {
	hub := sentry.GetHubFromContext(ctx)
	setting, err := ws.workspaceRepo.FetchPostmarkMailServerSettingById(ctx, workspaceId)

	if errors.Is(err, repository.ErrEmpty) {
		return models.PostmarkMailServerSetting{}, ErrPostmarkSettingNotFound
	}
	if errors.Is(err, repository.ErrQuery) {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, ErrPostmarkSetting
	}
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}
	return setting, nil
}

func (ws *WorkspaceService) PostmarkMailServerAddDomain(
	ctx context.Context, setting models.PostmarkMailServerSetting, domain string,
) (models.PostmarkMailServerSetting, bool, error) {
	var created bool
	var addedDomain postmark.DomainDetail
	var err error

	hub := sentry.GetHubFromContext(ctx)
	client := postmark.NewClient("", zyg.PostmarkAccountToken())

	// Checks if the DNS domain is already added.
	// If added we fetch the domain details from the Postmark API
	// otherwise, add the domain in Postmark
	if setting.DNSDomainId != nil {
		addedDomain, err = client.GetDomain(ctx, *setting.DNSDomainId)
		if err != nil {
			hub.CaptureException(err)
			return models.PostmarkMailServerSetting{}, created, err
		}
	} else {
		req := postmark.CreateDomainRequest{
			Name: domain,
		}
		addedDomain, err = client.CreateDomain(ctx, req)
		if err != nil {
			hub.CaptureException(err)
			return models.PostmarkMailServerSetting{}, created, err
		}
		created = true
	}

	// Pick the latest DKIM Host and TXT Value
	// As per the docs *Pending* should be the latest.
	var dkimHost, dkimTextValue string
	if addedDomain.DKIMPendingHost != "" {
		dkimHost = addedDomain.DKIMPendingHost
	} else {
		dkimHost = addedDomain.DKIMHost
	}
	if addedDomain.DKIMPendingTextValue != "" {
		dkimTextValue = addedDomain.DKIMPendingTextValue
	} else {
		dkimTextValue = addedDomain.DKIMTextValue
	}

	setting.HasDNS = true
	setting.IsDNSVerified = false
	setting.DNSDomainId = &addedDomain.ID
	setting.DKIMHost = &dkimHost
	setting.DKIMTextValue = &dkimTextValue
	setting.DKIMUpdateStatus = &addedDomain.DKIMUpdateStatus

	setting.ReturnPathDomain = &addedDomain.ReturnPathDomain
	setting.ReturnPathDomainCNAME = &addedDomain.ReturnPathDomainCNAMEValue
	setting.ReturnPathDomainVerified = addedDomain.ReturnPathDomainVerified

	setting, err = ws.workspaceRepo.SavePostmarkMailServerSetting(ctx, setting)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, created, err
	}
	return setting, created, nil
}

func (ws *WorkspaceService) PostmarkMailServerVerifyDomain(
	ctx context.Context, setting models.PostmarkMailServerSetting,
) (models.PostmarkMailServerSetting, error) {

	hub := sentry.GetHubFromContext(ctx)
	client := postmark.NewClient("", zyg.PostmarkAccountToken())

	now := time.Now().UTC()
	verifiedDKIM, err := client.VerifyDKIM(ctx, *setting.DNSDomainId)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}

	verifiedReturnPath, err := client.VerifyReturnPath(ctx, *setting.DNSDomainId)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}

	// Pick the latest DKIM Host and TXT Value
	// As per the docs *Pending* should be the latest.
	var dkimHost, dkimTextValue string
	if verifiedDKIM.DKIMPendingHost != "" {
		dkimHost = verifiedDKIM.DKIMPendingHost
	} else {
		dkimHost = verifiedDKIM.DKIMHost
	}
	if verifiedDKIM.DKIMPendingTextValue != "" {
		dkimTextValue = verifiedDKIM.DKIMPendingTextValue
	} else {
		dkimTextValue = verifiedDKIM.DKIMTextValue
	}

	setting.HasDNS = true
	setting.DNSVerifiedAt = &now
	setting.DKIMHost = &dkimHost
	setting.DKIMTextValue = &dkimTextValue
	setting.DKIMUpdateStatus = &verifiedDKIM.DKIMUpdateStatus

	setting.ReturnPathDomain = &verifiedReturnPath.ReturnPathDomain
	setting.ReturnPathDomainCNAME = &verifiedReturnPath.ReturnPathDomainCNAMEValue
	setting.ReturnPathDomainVerified = verifiedReturnPath.ReturnPathDomainVerified
	if setting.DNSHasVerified() {
		setting.IsDNSVerified = true
	} else {
		setting.IsDNSVerified = false
	}

	setting, err = ws.workspaceRepo.SavePostmarkMailServerSetting(ctx, setting)
	if err != nil {
		hub.CaptureException(err)
		return models.PostmarkMailServerSetting{}, err
	}
	return setting, nil
}

func (ws *WorkspaceService) PostmarkMailServerUpdate(
	ctx context.Context, setting models.PostmarkMailServerSetting, fields []string,
) (models.PostmarkMailServerSetting, error) {
	setting, err := ws.workspaceRepo.ModifyPostmarkMailServerSettingById(ctx, setting, fields)
	if err != nil {
		return models.PostmarkMailServerSetting{}, err
	}
	return setting, nil
}
