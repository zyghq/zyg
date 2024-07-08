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
}

type CustomerServicer interface {
	GetCustomerByExternalId(ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	GetCustomerByEmail(ctx context.Context, workspaceId string, email string) (models.Customer, error)
	GetCustomerByPhone(ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	GenerateCustomerToken(c models.Customer) (string, error)
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
	AddLabelToThread(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	ListThreadLabels(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	AddCustomerMessageToThread(ctx context.Context, th models.ThreadChat, c *models.Customer, msg string) (models.ThreadChatMessage, error)
	AddMemberMessageToThread(ctx context.Context, th models.ThreadChat, m *models.Member, msg string) (models.ThreadChatMessage, error)
	ListThreadMessages(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
	GenerateMemberThreadMetrics(ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error)
}

type AccountRepositorer interface {
	GetOrCreateByAuthUserId(ctx context.Context, account models.Account) (models.Account, bool, error)
	GetByAuthUserId(ctx context.Context, authUserId string) (models.Account, error)
	CreatePersonalAccessToken(ctx context.Context, ap models.AccountPAT) (models.AccountPAT, error)
	GetPatListByAccountId(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	GetAccountByToken(ctx context.Context, token string) (models.Account, error)
	GetPatByPatId(ctx context.Context, patId string) (models.AccountPAT, error)
	HardDeletePatByPatId(ctx context.Context, patId string) error
}

type WorkspaceRepositorer interface {
	CreateWorkspaceByAccount(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error)
	UpdateWorkspaceById(ctx context.Context, w models.Workspace) (models.Workspace, error)
	UpdateWorkspaceLabelById(ctx context.Context, workspaceId string, l models.Label) (models.Label, error)
	GetWorkspaceById(ctx context.Context, workspaceId string) (models.Workspace, error)
	GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	GetListByMemberAccountId(ctx context.Context, accountId string) ([]models.Workspace, error)
	GetOrCreateLabel(ctx context.Context, l models.Label) (models.Label, bool, error)
	GetWorkspaceLabelById(ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	GetLabelListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Label, error)
	AddMemberByWorkspaceId(ctx context.Context, workspaceId string, member models.Member) (models.Member, error)
}

type MemberRepositorer interface {
	GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (models.Member, error)
	GetListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Member, error)
	GetByWorkspaceMemberId(ctx context.Context, workspaceId string, memberId string) (models.Member, error)
}

type CustomerRepositorer interface {
	GetByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	GetWorkspaceCustomerByExtId(ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	GetWorkspaceCustomerByEmail(ctx context.Context, workspaceId string, email string) (models.Customer, error)
	GetWorkspaceCustomerByPhone(ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	GetOrCreateCustomerByExtId(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	GetOrCreateCustomerByEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	GetOrCreateCustomerByPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	GetListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.Customer, error)
}

type ThreadChatRepositorer interface {
	CreateThreadChat(ctx context.Context, th models.ThreadChat, msg string) (models.ThreadChat, models.ThreadChatMessage, error)
	GetByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (models.ThreadChat, error)
	UpdateThreadChatById(ctx context.Context, th models.ThreadChat, fields []string) (models.ThreadChat, error)
	GetListByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) ([]models.ThreadChatWithMessage, error)
	SetAssignee(ctx context.Context, threadChatId string, assigneeId string) (models.ThreadChat, error)
	SetReplied(ctx context.Context, threadChatId string, replied bool) (models.ThreadChat, error)
	GetListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	GetMemberAssignedListByWorkspaceId(ctx context.Context, workspaceId string, memberId string) ([]models.ThreadChatWithMessage, error)
	GetUnassignedListByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	GetLabelledListByWorkspaceId(ctx context.Context, worskapceId string, labelId string) ([]models.ThreadChatWithMessage, error)
	IsExistByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AddLabelToThread(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	GetLabelListByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	CreateCustomerThChatMessage(ctx context.Context, threadChatId string, customerId string, msg string) (models.ThreadChatMessage, error)
	CreateMemberThChatMessage(ctx context.Context, threadChatId string, memberId string, msg string) (models.ThreadChatMessage, error)
	GetMessageListByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
	StatusMetricsByWorkspaceId(ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	MemberAssigneeMetricsByWorkspaceId(ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	LabelMetricsByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
}
