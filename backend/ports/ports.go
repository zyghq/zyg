package ports

import (
	"context"

	"github.com/zyghq/zyg/models"
)

type AccountServicer interface {
	CreateAuthAccount(
		ctx context.Context, authUserId string, email string, name string, provider string,
	) (models.Account, bool, error)
	GeneratePersonalAccessToken(
		ctx context.Context, accountId string, name string, description string) (models.AccountPAT, error)
	ListPersonalAccessTokens(
		ctx context.Context, accountId string) ([]models.AccountPAT, error)
	GetPersonalAccessToken(
		ctx context.Context, patId string) (models.AccountPAT, error)
	DeletePersonalAccessToken(
		ctx context.Context, patId string) error
	CreateWorkspace(
		ctx context.Context, account models.Account, workspaceName string) (models.Workspace, error)
	GetAccountLinkedWorkspace(
		ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	ListAccountLinkedWorkspaces(
		ctx context.Context, accountId string) ([]models.Workspace, error)
}

type AuthServicer interface {
	AuthenticateUserAccount(
		ctx context.Context, authUserId string) (models.Account, error)
	AuthenticateWorkspaceMember(
		ctx context.Context, workspaceId string, accountId string) (models.Member, error)
	ValidatePersonalAccessToken(
		ctx context.Context, token string) (models.Account, error)
}

type CustomerAuthServicer interface {
	AuthenticateWorkspaceCustomer(
		ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error)
	GetWidgetLinkedSecretKey(
		ctx context.Context, widgetId string) (models.WorkspaceSecret, error)
}

type WorkspaceServicer interface {
	UpdateWorkspace(
		ctx context.Context, workspace models.Workspace) (models.Workspace, error)
	UpdateLabel(
		ctx context.Context, label models.Label) (models.Label, error)
	GetWorkspace(
		ctx context.Context, workspaceId string) (models.Workspace, error)
	CreateLabel(
		ctx context.Context, workspaceId string, name string, icon string) (models.Label, bool, error)
	GetLabel(
		ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	ListLabels(
		ctx context.Context, workspaceId string) ([]models.Label, error)
	GetAccountLinkedMember(
		ctx context.Context, workspaceId string, accountId string) (models.Member, error)
	ListMembers(
		ctx context.Context, workspaceId string) ([]models.Member, error)
	GetMember(
		ctx context.Context, workspaceId string, memberId string) (models.Member, error)
	GetSystemMember(
		ctx context.Context, workspaceId string) (models.Member, error)
	CreateNewSystemMember(
		ctx context.Context, workspaceId string) (models.Member, error)
	ListCustomers(
		ctx context.Context, workspaceId string) ([]models.Customer, error)
	CreateCustomerWithExternalId(
		ctx context.Context, workspaceId string, externalId string, isVerified bool, name string,
	) (models.Customer, bool, error)
	CreateCustomerWithEmail(
		ctx context.Context, workspaceId string, email string, isVerified bool, name string,
	) (models.Customer, bool, error)
	CreateCustomerWithPhone(
		ctx context.Context, workspaceId string, phone string, isVerified bool, name string,
	) (models.Customer, bool, error)
	CreateUnverifiedCustomer(
		ctx context.Context, workspaceId string, name string) (models.Customer, error)
	CreateWidget(
		ctx context.Context, workspaceId string, name string, configuration map[string]interface{},
	) (models.Widget, error)
	ListWidgets(ctx context.Context, workspaceId string) ([]models.Widget, error)
	GenerateWorkspaceSecret(
		ctx context.Context, workspaceId string, length int) (models.WorkspaceSecret, error)
	GetSecretKey(
		ctx context.Context, workspaceId string) (models.WorkspaceSecret, error)
	GetOrGenerateSecretKey(
		ctx context.Context, workspaceId string) (models.WorkspaceSecret, error)
	GetWidget(
		ctx context.Context, widgetId string) (models.Widget, error)
	GetCustomer(
		ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error)
	DoesEmailConflict(
		ctx context.Context, workspaceId string, email string) (bool, error)
	ValidateWidgetSession(
		ctx context.Context, sk string, widgetId string, sessionId string) (models.Customer, error)
	CreateWidgetSession(
		ctx context.Context, sk string, workspaceId string, widgetId string,
		sessionId string, name string) (models.Customer, bool, error)
}

type CustomerServicer interface {
	GenerateCustomerJwt(
		customer models.Customer, sk string) (string, error)
	VerifyExternalId(
		sk string, hash string, externalId string) bool
	VerifyEmail(
		sk string, hash string, email string) bool
	VerifyPhone(
		sk string, hash string, phone string) bool
	UpdateCustomer(
		ctx context.Context, customer models.Customer) (models.Customer, error)
	AddCustomerEmailIdentity(
		ctx context.Context, emailIdentity models.EmailIdentity) (models.EmailIdentity, error)
	HasProvidedEmailIdentity(ctx context.Context, customerId string) (bool, error)
}

type ThreadServicer interface {
	CreateNewInboundThreadChat(
		ctx context.Context, workspaceId string,
		customer models.Customer, createdBy models.MemberActor, message string,
	) (models.Thread, models.Chat, error)
	GetWorkspaceThread(
		ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error)
	UpdateThread(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)
	ListCustomerThreadChats(
		ctx context.Context, customerId string) ([]models.Thread, error)
	ListWorkspaceThreadChats(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListMemberThreadChats(
		ctx context.Context, memberId string) ([]models.Thread, error)
	ListUnassignedThreadChats(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListLabelledThreadChats(
		ctx context.Context, labelId string) ([]models.Thread, error)
	ThreadExistsInWorkspace(
		ctx context.Context, workspaceId string, threadId string) (bool, error)
	SetLabel(
		ctx context.Context, threadId string, labelId string, addedBy string) (models.ThreadLabel, bool, error)
	ListThreadLabels(
		ctx context.Context, threadChatId string) ([]models.ThreadLabel, error)
	CreateInboundChatMessage(
		ctx context.Context, thread models.Thread, message string) (models.Chat, error)
	CreateOutboundChatMessage(
		ctx context.Context, thread models.Thread, memberId string, message string) (models.Chat, error)
	ListThreadChatMessages(
		ctx context.Context, threadId string) ([]models.Chat, error)
	GenerateMemberThreadMetrics(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error)
	RemoveThreadLabel(
		ctx context.Context, threadId string, labelId string) error
}

type AccountRepositorer interface {
	UpsertByAuthUserId(
		ctx context.Context, account models.Account) (models.Account, bool, error)
	FetchByAuthUserId(
		ctx context.Context, authUserId string) (models.Account, error)
	InsertPersonalAccessToken(
		ctx context.Context, pat models.AccountPAT) (models.AccountPAT, error)
	FetchPatsByAccountId(
		ctx context.Context, accountId string) ([]models.AccountPAT, error)
	LookupByToken(
		ctx context.Context, token string) (models.Account, error)
	FetchPatById(
		ctx context.Context, patId string) (models.AccountPAT, error)
	DeletePatById(
		ctx context.Context, patId string) error
}

type WorkspaceRepositorer interface {
	InsertWorkspaceWithMembers(
		ctx context.Context, workspace models.Workspace, members []models.Member) (models.Workspace, error)
	ModifyWorkspaceById(
		ctx context.Context, workspace models.Workspace) (models.Workspace, error)
	ModifyLabelById(
		ctx context.Context, label models.Label) (models.Label, error)
	FetchByWorkspaceId(
		ctx context.Context, workspaceId string) (models.Workspace, error)
	LookupWorkspaceByAccountId(
		ctx context.Context, workspaceId string, accountId string) (models.Workspace, error)
	FetchWorkspacesByAccountId(
		ctx context.Context, accountId string) ([]models.Workspace, error)
	InsertLabelByName(
		ctx context.Context, label models.Label) (models.Label, bool, error)
	LookupWorkspaceLabelById(
		ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	FetchLabelsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Label, error)
	InsertWidget(
		ctx context.Context, widget models.Widget) (models.Widget, error)
	FetchWidgetsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Widget, error)
	InsertWorkspaceSecret(
		ctx context.Context, workspaceId string, sk string) (models.WorkspaceSecret, error)
	FetchSecretKeyByWorkspaceId(
		ctx context.Context, workspaceId string) (models.WorkspaceSecret, error)
	LookupWidgetById(
		ctx context.Context, widgetId string) (models.Widget, error)
	LookupWidgetSessionById(
		ctx context.Context, widgetId string, sessionId string) (models.WidgetSession, error)
	UpsertWidgetSessionById(
		ctx context.Context, session models.WidgetSession) (models.WidgetSession, bool, error)
	InsertSystemMember(
		ctx context.Context, member models.Member) (models.Member, error)
	LookupSystemMemberByOldest(
		ctx context.Context, workspaceId string) (models.Member, error)
}

type MemberRepositorer interface {
	LookupByWorkspaceAccountId(
		ctx context.Context, workspaceId string, accountId string) (models.Member, error)
	FetchMembersByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Member, error)
	FetchByWorkspaceMemberId(
		ctx context.Context, workspaceId string, memberId string) (models.Member, error)
}

type CustomerRepositorer interface {
	LookupWorkspaceCustomerById(
		ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error)
	UpsertCustomerByExtId(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerByEmail(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerByPhone(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerById(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	FetchCustomersByWorkspaceId(
		ctx context.Context, workspaceId string, role *string) ([]models.Customer, error)
	LookupSecretKeyByWidgetId(
		ctx context.Context, widgetId string) (models.WorkspaceSecret, error)
	ModifyCustomerById(
		ctx context.Context, customer models.Customer) (models.Customer, error)
	CheckEmailExists(
		ctx context.Context, workspaceId string, email string) (bool, error)
	InsertEmailIdentity(ctx context.Context, identity models.EmailIdentity) (models.EmailIdentity, error)
	EmailIdentityExists(ctx context.Context, customerId string) (bool, error)
}

type ThreadRepositorer interface {
	InsertInboundThreadChat(
		ctx context.Context, thread models.Thread, chat models.Chat) (models.Thread, models.Chat, error)
	LookupByWorkspaceThreadId(
		ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error)
	ModifyThreadById(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)
	FetchThreadsByCustomerId(
		ctx context.Context, customerId string, channel *string) ([]models.Thread, error)
	FetchThreadsByWorkspaceId(
		ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error)
	FetchThreadsByAssignedMemberId(
		ctx context.Context, memberId string, channel *string, role *string) ([]models.Thread, error)
	FetchThreadsByMemberUnassigned(
		ctx context.Context, workspaceId string, channel *string, role *string) ([]models.Thread, error)
	FetchThreadsByLabelId(
		ctx context.Context, labelId string, channel *string, role *string) ([]models.Thread, error)
	CheckThreadInWorkspaceExists(
		ctx context.Context, workspaceId string, threadId string) (bool, error)
	SetThreadLabel(
		ctx context.Context, thl models.ThreadLabel) (models.ThreadLabel, bool, error)
	FetchAttachedLabelsByThreadId(
		ctx context.Context, threadId string) ([]models.ThreadLabel, error)
	InsertCustomerChat(
		ctx context.Context, thread models.Thread, inboundMessage models.InboundMessage, chat models.Chat,
	) (models.Chat, error)
	InsertMemberChat(
		ctx context.Context, thread models.Thread, outboundMessage models.OutboundMessage, chat models.Chat,
	) (models.Chat, error)
	FetchThChatMessagesByThreadId(
		ctx context.Context, threadId string) ([]models.Chat, error)
	ComputeStatusMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	ComputeAssigneeMetricsByMember(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	ComputeLabelMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
	DeleteThreadLabelById(
		ctx context.Context, threadId string, labelId string) error
}
