package ports

import (
	"context"

	"github.com/zyghq/zyg/models"
)

type AccountServicer interface {
	InitiateAccount(ctx context.Context, a models.Account) (models.Account, bool, error)
	AuthUser(ctx context.Context, authUserId string) (models.Account, error)
	IssuePersonalAccessToken(ctx context.Context, ap models.AccountPAT) (models.AccountPAT, error)
	UserPats(ctx context.Context, accountId string) ([]models.AccountPAT, error)
	PatAccount(ctx context.Context, token string) (models.Account, error)
	UserPat(ctx context.Context, patId string) (models.AccountPAT, error)
	HardDeletePat(ctx context.Context, patId string) error
}

type AuthServicer interface {
	CheckAuthUser(ctx context.Context, authUserId string) (models.Account, error)
	CheckPatAccount(ctx context.Context, token string) (models.Account, error)
}

type CustomerAuthServicer interface {
	WorkspaceCustomer(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
}

type WorkspaceServicer interface {
	CreateAccountWorkspace(ctx context.Context, a models.Account, w models.Workspace) (models.Workspace, error)
	UpdateWorkspace(ctx context.Context, w models.Workspace) (models.Workspace, error)
	UpdateWorkspaceLabel(ctx context.Context, workspaceId string, label models.Label) (models.Label, error)
	GetWorkspace(ctx context.Context, workspaceId string) (models.Workspace, error)
	MemberWorkspace(ctx context.Context, accountId string, workspaceId string) (models.Workspace, error)
	MemberWorkspaces(ctx context.Context, accountId string) ([]models.Workspace, error)
	InitWorkspaceLabel(ctx context.Context, label models.Label) (models.Label, bool, error)
	WorkspaceLabel(ctx context.Context, workspaceId string, labelId string) (models.Label, error)
	WorkspaceLabels(ctx context.Context, workspaceId string) ([]models.Label, error)
	WorkspaceUserMember(ctx context.Context, accountId string, workspaceId string) (models.Member, error) // TODO: perhaps rename to WorkspaceAccountMember?
	AddMember(ctx context.Context, workspace models.Workspace, member models.Member) (models.Member, error)
	WorkspaceMembers(ctx context.Context, workspaceId string) ([]models.Member, error)
	WorkspaceMember(ctx context.Context, workspaceId string, memberId string) (models.Member, error)
	WorkspaceCustomers(ctx context.Context, workspaceId string) ([]models.Customer, error)
	InitWorkspaceCustomerWithExternalId(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	InitWorkspaceCustomerWithEmail(ctx context.Context, c models.Customer) (models.Customer, bool, error)
	InitWorkspaceCustomerWithPhone(ctx context.Context, c models.Customer) (models.Customer, bool, error)
}

type CustomerServicer interface {
	WorkspaceCustomer(ctx context.Context, workspaceId string, customerId string) (models.Customer, error)
	WorkspaceCustomerWithExternalId(ctx context.Context, workspaceId string, externalId string) (models.Customer, error)
	WorkspaceCustomerWithEmail(ctx context.Context, workspaceId string, email string) (models.Customer, error)
	WorkspaceCustomerWithPhone(ctx context.Context, workspaceId string, phone string) (models.Customer, error)
	IssueCustomerJwt(c models.Customer) (string, error)
}

type ThreadChatServicer interface {
	CreateCustomerThread(ctx context.Context, th models.ThreadChat, msg string) (models.ThreadChat, models.ThreadChatMessage, error)
	WorkspaceThread(ctx context.Context, workspaceId string, threadChatId string) (models.ThreadChat, error)
	UpdateThreadChat(ctx context.Context, th models.ThreadChat, fields []string) (models.ThreadChat, error)
	WorkspaceCustomerThreadChats(ctx context.Context, workspaceId string, customerId string) ([]models.ThreadChatWithMessage, error)
	AssignMember(ctx context.Context, threadChatId string, assigneeId string) (models.ThreadChat, error)
	MarkReplied(ctx context.Context, threadChatId string, replied bool) (models.ThreadChat, error)
	WorkspaceThreads(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	WorkspaceMemberAssignedThreadList(ctx context.Context, workspaceId string, memberId string) ([]models.ThreadChatWithMessage, error)
	WorkspaceUnassignedThreadList(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error)
	WorkspaceLabelledThreadList(ctx context.Context, workspaceId string, labelId string) ([]models.ThreadChatWithMessage, error)
	ExistInWorkspace(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AddLabel(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	ThreadLabels(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	CreateCustomerMessage(ctx context.Context, th models.ThreadChat, c *models.Customer, msg string) (models.ThreadChatMessage, error)
	CreateMemberMessage(ctx context.Context, th models.ThreadChat, m *models.Member, msg string) (models.ThreadChatMessage, error)
	ThreadChatMessages(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
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
	AddLabel(ctx context.Context, thl models.ThreadChatLabel) (models.ThreadChatLabel, bool, error)
	GetLabelListByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error)
	CreateCustomerThChatMessage(ctx context.Context, threadChatId string, customerId string, msg string) (models.ThreadChatMessage, error)
	CreateMemberThChatMessage(ctx context.Context, threadChatId string, memberId string, msg string) (models.ThreadChatMessage, error)
	GetMessageListByThreadChatId(ctx context.Context, threadChatId string) ([]models.ThreadChatMessage, error)
	StatusMetricsByWorkspaceId(ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	MemberAssigneeMetricsByWorkspaceId(ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	LabelMetricsByWorkspaceId(ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
}
