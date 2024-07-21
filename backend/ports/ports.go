package ports

import (
	"context"

	"github.com/zyghq/zyg/models"
)

type AccountServicer interface {
	InitiateAccount(ctx context.Context, a models.Account) (models.Account, bool, error)
	GeneratePersonalAccessToken(ctx context.Context, ap models.AccountPAT) (models.AccountPAT, error)
	GetPersonalAccessTokens(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	GetPersonalAccessToken(ctx context.Context, patId string) (models.AccountPAT, error)
	DeletePersonalAccessToken(ctx context.Context, patId string) error
}

type AuthServicer interface {
	AuthenticateUser(ctx context.Context, authUserId string) (models.Account, error)
	ValidatePersonalAccessToken(ctx context.Context, token string) (models.Account, error)
}

type CustomerAuthServicer interface {
	ValidateWorkspaceCustomer(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	GetWidgetLinkedSecretKey(ctx context.Context, widgetId string) (models.SecretKey, error)
}

type WorkspaceServicer interface {
	CreateWorkspace(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error)
	UpdateWorkspace(ctx context.Context, w models.Workspace) (models.Workspace, error)
	SetWorkspaceLabel(ctx context.Context, workspaceId string, label models.Label) (models.Label, error)
	GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error)
	GetMemberWorkspace(ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	ListMemberWorkspaces(ctx context.Context, accountId string) ([]models.Workspace, error)
	CreateLabel(ctx context.Context, label models.Label) (models.Label, bool, error)
	GetWorkspaceLabel(ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	ListWorkspaceLabels(ctx context.Context, workspaceId string) ([]models.Label, error)
	GetWorkspaceMember(ctx context.Context, accountId string, workspaceId string) (models.Member, error)
	AddMember(ctx context.Context, workspace models.Workspace, member models.Member) (models.Member, error)
	ListWorkspaceMembers(ctx context.Context, workspaceId string) ([]models.Member, error)
	GetWorkspaceMemberById(ctx context.Context, workspaceId string, memberId string) (models.Member, error)
	ListWorkspaceCustomers(ctx context.Context, workspaceId string) ([]models.Customer, error)
	CreateCustomerByExternalId(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	CreateWorkspaceCustomerWithEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	CreateWorkspaceCustomerWithPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	CreateWidget(ctx context.Context, workspaceId string, widget models.Widget) (models.Widget, error)
	ListWidgets(ctx context.Context, workspaceId string) ([]models.Widget, error)
	GenerateSecretKey(ctx context.Context, workspaceId string, length int) (models.SecretKey, error)
	GetWorkspaceSecretKey(ctx context.Context, workspaceId string) (models.SecretKey, error)
	GetWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error)
}

type CustomerServicer interface {
	GetCustomerByExternalId(ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	GetCustomerByEmail(ctx context.Context, workspaceId string, email string) (models.Customer, error)
	GetCustomerByPhone(ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	GenerateCustomerToken(c models.Customer, sk string) (string, error)
	VerifyExternalId(sk string, hash string, externalId string) bool
	VerifyEmail(sk string, hash string, email string) bool
	VerifyPhone(sk string, hash string, phone string) bool
}

type ThreadChatServicer interface {
	CreateThreadWithMessage(ctx context.Context, th models.ThreadChat, msg string) (models.ThreadChat, models.ThreadChatMessage, error)
	GetThread(ctx context.Context, workspaceId string, threadChatId string) (models.ThreadChat, error)
	UpdateThread(ctx context.Context, th models.ThreadChat, fields []string) (models.ThreadChat, error)
	ListCustomerThreads(ctx context.Context, workspaceId string, customerId string) ([]models.ThreadChatWithMessage, error)
	AssignMemberToThread(ctx context.Context, threadChatId string, assigneeId string) (models.ThreadChat, error)
	SetThreadReplyStatus(ctx context.Context, threadChatId string, replied bool) (models.ThreadChat, error)
	ListWorkspaceThreads(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	ListMemberAssignedThreads(ctx context.Context, workspaceId string, memberId string) ([]models.ThreadChatWithMessage, error)
	ListUnassignedThreads(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	ListLabelledThreads(ctx context.Context, workspaceId string, labelId string) ([]models.ThreadChatWithMessage, error)
	ThreadExistsInWorkspace(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AttachLabelToThread(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	ListThreadLabels(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	AddCustomerMessageToThread(ctx context.Context, th models.ThreadChat, c *models.Customer, msg string) (models.ThreadChatMessage, error)
	AddMemberMessageToThread(ctx context.Context, th models.ThreadChat, m *models.Member, msg string) (models.ThreadChatMessage, error)
	ListThreadMessages(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
	GenerateMemberThreadMetrics(ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error)
}

type AccountRepositorer interface {
	UpsertAccountByAuthId(ctx context.Context, account models.Account) (models.Account, bool, error)
	FetchAccountByAuthId(ctx context.Context, authUserId string) (models.Account, error)
	InsertPersonalAccessToken(ctx context.Context, ap models.AccountPAT) (models.AccountPAT, error)
	RetrievePatsByAccountId(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	LookupAccountByToken(ctx context.Context, token string) (models.Account, error)
	FetchPatByPatId(ctx context.Context, patId string) (models.AccountPAT, error)
	PermanentlyRemovePatByPatId(ctx context.Context, patId string) error
}

type WorkspaceRepositorer interface {
	InsertWorkspaceForAccount(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error)
	ModifyWorkspaceById(ctx context.Context, w models.Workspace) (models.Workspace, error)
	AlterWorkspaceLabelById(ctx context.Context, workspaceId string, l models.Label) (models.Label, error)
	FetchWorkspaceById(ctx context.Context, workspaceId string) (models.Workspace, error)
	RetrieveByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	FetchWorkspacesByMemberAccountId(ctx context.Context, accountId string) ([]models.Workspace, error)
	UpsertLabel(ctx context.Context, l models.Label) (models.Label, bool, error)
	LookupWorkspaceLabelById(ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	RetrieveLabelsByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Label, error)
	InsertMemberIntoWorkspace(ctx context.Context, workspaceId string, member models.Member) (models.Member, error)
	InsertWidgetIntoWorkspace(ctx context.Context, workspaceId string, widget models.Widget) (models.Widget, error)
	RetrieveWidgetsByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Widget, error)
	InsertSecretKeyIntoWorkspace(ctx context.Context, workspaceId string, sk string) (models.SecretKey, error)
	FetchSecretKeyByWorkspaceId(ctx context.Context, workspaceId string) (models.SecretKey, error)
	LookupWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error)
}

type MemberRepositorer interface {
	LookupByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (models.Member, error)
	RetrieveMembersByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Member, error)
	FetchByWorkspaceMemberId(ctx context.Context, workspaceId string, memberId string) (models.Member, error)
}

type CustomerRepositorer interface {
	LookupByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	FetchWorkspaceCustomerByExtId(ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	RetrieveWorkspaceCustomerByEmail(ctx context.Context, workspaceId string, email string) (models.Customer, error)
	LookupWorkspaceCustomerByPhone(ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	UpsertCustomerByExtId(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	UpsertCustomerByEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	UpsertCustomerByPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	FetchCustomersByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Customer, error)
	LookupSecretKeyByWidgetId(ctx context.Context, widgetId string) (models.SecretKey, error)
}

type ThreadChatRepositorer interface {
	InsertThreadChat(ctx context.Context, th models.ThreadChat, msg string) (models.ThreadChat, models.ThreadChatMessage, error)
	LookupByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (models.ThreadChat, error)
	ModifyThreadChatById(ctx context.Context, th models.ThreadChat, fields []string) (models.ThreadChat, error)
	FetchByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) ([]models.ThreadChatWithMessage, error)
	UpdateAssignee(ctx context.Context, threadChatId string, assigneeId string) (models.ThreadChat, error)
	UpdateRepliedStatus(ctx context.Context, threadChatId string, replied bool) (models.ThreadChat, error)
	RetrieveByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	FetchAssignedThreadsByMember(ctx context.Context, workspaceId string, memberId string) ([]models.ThreadChatWithMessage, error)
	RetrieveUnassignedThreads(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	FetchThreadsByLabel(ctx context.Context, worskapceId string, labelId string) ([]models.ThreadChatWithMessage, error)
	CheckExistenceByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AttachLabelToThread(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	RetrieveLabelsByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	InsertCustomerMessage(ctx context.Context, threadChatId string, customerId string, msg string) (models.ThreadChatMessage, error)
	InsertMemberMessage(ctx context.Context, threadChatId string, memberId string, msg string) (models.ThreadChatMessage, error)
	FetchMessagesByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
	ComputeStatusMetricsByWorkspaceId(ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	CalculateAssigneeMetricsByMember(ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	ComputeLabelMetricsByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
}
