package ports

import (
	"context"

	"github.com/zyghq/zyg/models"
)

type AccountServicer interface {
	CreateAuthAccount(
		ctx context.Context, authUserId string, email string, name string, provider string) (models.Account, bool, error)
	GeneratePersonalAccessToken(
		ctx context.Context, accountId string, name string, description string) (models.AccountPAT, error)
	GetPersonalAccessTokens(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	GetPersonalAccessToken(
		ctx context.Context, patId string) (models.AccountPAT, error)
	DeletePersonalAccessToken(ctx context.Context, patId string) error
}

type AuthServicer interface {
	AuthenticateUser(ctx context.Context, authUserId string) (models.Account, error)
	ValidatePersonalAccessToken(ctx context.Context, token string) (models.Account, error)
}

type CustomerAuthServicer interface {
	GetWorkspaceCustomerIgnoreRole(
		ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	GetWidgetLinkedSecretKey(
		ctx context.Context, widgetId string) (models.SecretKey, error)
}

type WorkspaceServicer interface {
	CreateWorkspace(
		ctx context.Context, accountId string, memberName string, workspaceName string) (models.Workspace, error)
	UpdateWorkspace(
		ctx context.Context, workspace models.Workspace) (models.Workspace, error)
	UpdateWorkspaceLabel(
		ctx context.Context, workspaceId string, label models.Label) (models.Label, error)
	GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error)
	GetLinkedWorkspaceMember(
		ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	ListAccountWorkspaces(
		ctx context.Context, accountId string) ([]models.Workspace, error)
	CreateLabel(
		ctx context.Context, workspaceId string, name string, icon string) (models.Label, bool, error)
	GetWorkspaceLabel(
		ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	ListWorkspaceLabels(
		ctx context.Context, workspaceId string) ([]models.Label, error)
	GetWorkspaceAccountMember(
		ctx context.Context, accountId string, workspaceId string) (models.Member, error)
	AddMember(ctx context.Context, workspace models.Workspace, member models.Member) (models.Member, error)
	ListWorkspaceMembers(
		ctx context.Context, workspaceId string) ([]models.Member, error)
	GetWorkspaceMemberById(ctx context.Context, workspaceId string, memberId string) (models.Member, error)
	ListWorkspaceCustomers(
		ctx context.Context, workspaceId string) ([]models.Customer, error)
	CreateCustomerWithExternalId(
		ctx context.Context, workspaceId string, externalId string, isVerified bool, name string) (models.Customer, bool, error)
	CreateCustomerWithEmail(
		ctx context.Context, workspaceId string, email string, isVerified bool, name string) (models.Customer, bool, error)
	CreateCustomerWithPhone(
		ctx context.Context, workspaceId string, phone string, isVerified bool, name string) (models.Customer, bool, error)
	CreateAnonymousCustomer(
		ctx context.Context, workspaceId string, anonId string, isVerified bool, name string) (models.Customer, bool, error)
	CreateWidget(
		ctx context.Context, workspaceId string, name string, configuration map[string]interface{}) (models.Widget, error)
	ListWidgets(ctx context.Context, workspaceId string) ([]models.Widget, error)
	GenerateSecretKey(
		ctx context.Context, workspaceId string, length int) (models.SecretKey, error)
	GetWorkspaceSecretKey(
		ctx context.Context, workspaceId string) (models.SecretKey, error)
	GetWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error)
}

type CustomerServicer interface {
	GetWorkspaceCustomerByExtId(
		ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	GetWorkspaceCustomerByEmail(
		ctx context.Context, workspaceId string, email string) (models.Customer, error)
	GetWorkspaceCustomerByPhone(
		ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	GenerateCustomerToken(
		customer models.Customer, sk string) (string, error)
	VerifyExternalId(sk string, hash string, externalId string) bool
	VerifyEmail(sk string, hash string, email string) bool
	VerifyPhone(sk string, hash string, phone string) bool
	GetWorkspaceCustomerById(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	UpdateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)
}

type ThreadChatServicer interface {
	CreateThreadInAppChat(
		ctx context.Context, workspaceId string, customerId string, message string,
	) (models.Thread, models.Chat, error)
	GetWorkspaceThread(
		ctx context.Context, workspaceId string, threadId string) (models.Thread, error)
	UpdateThread(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)
	ListCustomerThreadChats(ctx context.Context, workspaceId string, customerId string) ([]models.Thread, error)
	AssignMemberToThread(ctx context.Context, threadId string, assigneeId string) (models.Thread, error)
	SetThreadReplyStatus(ctx context.Context, threadId string, replied bool) (models.Thread, error)
	ListWorkspaceThreadChats(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListMemberAssignedThreadChats(
		ctx context.Context, memberId string) ([]models.Thread, error)
	ListUnassignedThreadChats(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListLabelledThreadChats(
		ctx context.Context, labelId string) ([]models.Thread, error)
	ThreadExistsInWorkspace(
		ctx context.Context, workspaceId string, threadId string) (bool, error)
	AttachLabelToThread(
		ctx context.Context, threadId string, labelId string, addedBy string) (models.ThreadLabel, bool, error)
	ListThreadLabels(ctx context.Context, threadChatId string) ([]models.ThreadLabel, error)
	AddCustomerMessageToThread(ctx context.Context, threadId string, customerId string, message string) (models.Chat, error)
	AddMemberMessageToThreadChat(
		ctx context.Context, threadId string, memberId string, message string) (models.Chat, error)
	ListThreadChatMessages(
		ctx context.Context, threadId string) ([]models.Chat, error)
	GenerateMemberThreadMetrics(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error)
}

type AccountRepositorer interface {
	UpsertByAuthUserId(ctx context.Context, account models.Account) (models.Account, bool, error)
	FetchAccountByAuthId(ctx context.Context, authUserId string) (models.Account, error)
	InsertPersonalAccessToken(ctx context.Context, pat models.AccountPAT) (models.AccountPAT, error)
	RetrievePatsByAccountId(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	LookupAccountByToken(ctx context.Context, token string) (models.Account, error)
	FetchPatById(ctx context.Context, patId string) (models.AccountPAT, error)
	DeletePatById(ctx context.Context, patId string) error
}

type WorkspaceRepositorer interface {
	InsertWorkspaceByAccountId(
		ctx context.Context, memberName string, workspace models.Workspace) (models.Workspace, error)
	ModifyWorkspaceById(ctx context.Context, w models.Workspace) (models.Workspace, error)
	ModifyWorkspaceLabelById(
		ctx context.Context, workspaceId string, l models.Label) (models.Label, error)
	FetchWorkspaceById(ctx context.Context, workspaceId string) (models.Workspace, error)
	LookupLinkedWorkspaceByAccountId(
		ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	FetchLinkedWorkspacesByAccountId(ctx context.Context, accountId string) ([]models.Workspace, error)
	InsertLabelByName(
		ctx context.Context, label models.Label) (models.Label, bool, error)
	LookupWorkspaceLabelById(
		ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	RetrieveLabelsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Label, error)
	InsertMemberIntoWorkspace(ctx context.Context, workspaceId string, member models.Member) (models.Member, error)
	InsertWidgetIntoWorkspace(
		ctx context.Context, workspaceId string, widget models.Widget) (models.Widget, error)
	RetrieveWidgetsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Widget, error)
	InsertSecretKeyIntoWorkspace(
		ctx context.Context, workspaceId string, sk string) (models.SecretKey, error)
	FetchSecretKeyByWorkspaceId(
		ctx context.Context, workspaceId string) (models.SecretKey, error)
	LookupWorkspaceWidget(ctx context.Context, widgetId string) (models.Widget, error)
}

type MemberRepositorer interface {
	LookupByAccountWorkspaceId(
		ctx context.Context, accountId string, workspaceId string) (models.Member, error)
	RetrieveMembersByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Member, error)
	FetchByWorkspaceMemberId(
		ctx context.Context, workspaceId string, memberId string) (models.Member, error)
}

type CustomerRepositorer interface {
	LookupByWorkspaceCustomerId(
		ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	LookupWorkspaceCustomerWithoutRoleById(
		ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	UpsertCustomerByExtId(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerByEmail(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerByPhone(
		ctx context.Context, customer models.Customer) (models.Customer, bool, error)
	UpsertCustomerByAnonId(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	FetchCustomersByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Customer, error)
	LookupSecretKeyByWidgetId(
		ctx context.Context, widgetId string) (models.SecretKey, error)
	ModifyCustomerById(ctx context.Context, c models.Customer) (models.Customer, error)
	LookupWorkspaceCustomerByExtId(
		ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	LookupWorkspaceCustomerByEmail(
		ctx context.Context, workspaceId string, email string) (models.Customer, error)
	LookupWorkspaceCustomerByPhone(
		ctx context.Context, workspaceId string, phone string) (models.Customer, error)
}

// todo: rename to ThreadRepository
type ThreadChatRepositorer interface {
	InsertInAppThreadChat(ctx context.Context, th models.Thread, chat models.Chat) (models.Thread, models.Chat, error) // done
	LookupByWorkspaceThreadId(
		ctx context.Context, workspaceId string, threadId string) (models.Thread, error)
	ModifyThreadById(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)
	// todo: rename FetchThChatsByCustomerId remove workspaceId
	RetrieveWorkspaceThChatsByCustomerId(
		ctx context.Context, workspaceId string, customerId string,
	) ([]models.Thread, error)
	UpdateAssignee(ctx context.Context, threadId string, assigneeId string) (models.Thread, error)
	UpdateRepliedStatus(ctx context.Context, threadId string, replied bool) (models.Thread, error)
	FetchThChatsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	FetchAssignedThChatsByMemberId(
		ctx context.Context, memberId string) ([]models.Thread, error)
	FetchUnassignedThChatsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	FetchThChatsByLabelId(
		ctx context.Context, labelId string) ([]models.Thread, error)
	CheckWorkspaceExistenceByThreadId(
		ctx context.Context, workspaceId string, threadId string) (bool, error)
	SetLabelToThread(
		ctx context.Context, thl models.ThreadLabel) (models.ThreadLabel, bool, error)
	RetrieveLabelsByThreadId(ctx context.Context, threadId string) ([]models.ThreadLabel, error)
	InsertCustomerMessage(ctx context.Context, chat models.Chat) (models.Chat, error)
	InsertThChatMemberMessage(
		ctx context.Context, chat models.Chat) (models.Chat, error)
	FetchThChatMessagesByThreadId(
		ctx context.Context, threadId string) ([]models.Chat, error)
	ComputeStatusMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	ComputeAssigneeMetricsByMember(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	ComputeLabelMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
}
