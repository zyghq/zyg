package ports

import (
	"context"

	"github.com/zyghq/zyg/internal/domain"
)

type AccountServicer interface {
	InitiateAccount(ctx context.Context, a domain.Account) (domain.Account, bool, error)
	GetAuthUser(ctx context.Context, authUserId string) (domain.Account, error)
	IssuePersonalAccessToken(ctx context.Context, ap domain.AccountPAT) (domain.AccountPAT, error)
	GetUserPatList(ctx context.Context, accountId string) ([]domain.AccountPAT, error)
	GetPatAccount(ctx context.Context, token string) (domain.Account, error)
}

type AuthServicer interface {
	GetAuthUser(ctx context.Context, authUserId string) (domain.Account, error)
	GetPatAccount(ctx context.Context, token string) (domain.Account, error)
}

type CustomerAuthServicer interface {
	GetWorkspaceCustomer(ctx context.Context, workspaceId string, customerId string) (domain.Customer, error)
}

type WorkspaceServicer interface {
	CreateWorkspace(ctx context.Context, w domain.Workspace) (domain.Workspace, error)
	GetWorkspace(ctx context.Context, workspaceId string) (domain.Workspace, error)
	GetUserWorkspace(ctx context.Context, accountId string, workspaceId string) (domain.Workspace, error)
	GetUserWorkspaceList(ctx context.Context, accountId string) ([]domain.Workspace, error)
	InitWorkspaceLabel(ctx context.Context, label domain.Label) (domain.Label, bool, error)
	GetWorkspaceMember(ctx context.Context, accountId string, workspaceId string) (domain.Member, error)
	InitWorkspaceCustomerWithExternalId(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
	InitWorkspaceCustomerWithEmail(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
	InitWorkspaceCustomerWithPhone(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
}

type CustomerServicer interface {
	GetWorkspaceCustomer(ctx context.Context, workspaceId string, customerId string) (domain.Customer, error)
	GetWorkspaceCustomerWithExternalId(ctx context.Context, workspaceId string, externalId string) (domain.Customer, error)
	GetWorkspaceCustomerWithEmail(ctx context.Context, workspaceId string, email string) (domain.Customer, error)
	GetWorkspaceCustomerWithPhone(ctx context.Context, workspaceId string, phone string) (domain.Customer, error)
	IssueJwt(c domain.Customer) (string, error)
}

type ThreadChatServicer interface {
	CreateCustomerThread(ctx context.Context, th domain.ThreadChat, msg string) (domain.ThreadChat, domain.ThreadChatMessage, error)
	GetWorkspaceThread(ctx context.Context, workspaceId string, threadChatId string) (domain.ThreadChat, error)
	GetWorkspaceCustomerList(ctx context.Context, workspaceId string, customerId string) ([]domain.ThreadChatWithMessage, error)
	AssignMember(ctx context.Context, threadChatId string, assigneeId string) (domain.ThreadChat, error)
	MarkReplied(ctx context.Context, threadChatId string, replied bool) (domain.ThreadChat, error)
	GetWorkspaceList(ctx context.Context, workspaceId string) ([]domain.ThreadChatWithMessage, error)
	ExistInWorkspace(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AddLabel(ctx context.Context, thl domain.ThreadChatLabel) (domain.ThreadChatLabel, bool, error)
	GetLabelList(ctx context.Context, threadChatId string) ([]domain.ThreadChatLabelled, error)
	CreateCustomerMessage(ctx context.Context, th domain.ThreadChat, c domain.Customer, msg string) (domain.ThreadChatMessage, error)
	CreateMemberMessage(ctx context.Context, th domain.ThreadChat, m domain.Member, msg string) (domain.ThreadChatMessage, error)
	GetMessageList(ctx context.Context, threadChatId string) ([]domain.ThreadChatMessage, error)
}

type AccountRepositorer interface {
	GetOrCreateByAuthUserId(ctx context.Context, account domain.Account) (domain.Account, bool, error)
	GetByAuthUserId(ctx context.Context, authUserId string) (domain.Account, error)
	CreatePersonalAccessToken(ctx context.Context, ap domain.AccountPAT) (domain.AccountPAT, error)
	GetPatListByAccountId(ctx context.Context, accountId string) ([]domain.AccountPAT, error)
	GetAccountByToken(ctx context.Context, token string) (domain.Account, error)
}

type WorkspaceRepositorer interface {
	CreateWorkspace(ctx context.Context, w domain.Workspace) (domain.Workspace, error)
	GetWorkspaceById(ctx context.Context, workspaceId string) (domain.Workspace, error)
	GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (domain.Workspace, error)
	GetListByAccountId(ctx context.Context, accountId string) ([]domain.Workspace, error)
	GetOrCreateLabel(ctx context.Context, l domain.Label) (domain.Label, bool, error)
}

type MemberRepositorer interface {
	GetByAccountWorkspaceId(ctx context.Context, accountId string, workspaceId string) (domain.Member, error)
}

type CustomerRepositorer interface {
	GetByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) (domain.Customer, error)
	GetWorkspaceCustomerByExtId(ctx context.Context, workspaceId string, externalId string) (domain.Customer, error)
	GetWorkspaceCustomerByEmail(ctx context.Context, workspaceId string, email string) (domain.Customer, error)
	GetWorkspaceCustomerByPhone(ctx context.Context, workspaceId string, phone string) (domain.Customer, error)
	GetOrCreateCustomerByExtId(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
	GetOrCreateCustomerByEmail(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
	GetOrCreateCustomerByPhone(ctx context.Context, c domain.Customer) (domain.Customer, bool, error)
}

type ThreadChatRepositorer interface {
	CreateThreadChat(ctx context.Context, th domain.ThreadChat, msg string) (domain.ThreadChat, domain.ThreadChatMessage, error)
	GetByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (domain.ThreadChat, error)
	GetListByWorkspaceCustomerId(ctx context.Context, workspaceId string, customerId string) ([]domain.ThreadChatWithMessage, error)
	SetAssignee(ctx context.Context, threadChatId string, assigneeId string) (domain.ThreadChat, error)
	SetReplied(ctx context.Context, threadChatId string, replied bool) (domain.ThreadChat, error)
	GetListByWorkspaceId(ctx context.Context, workspaceId string) ([]domain.ThreadChatWithMessage, error)
	IsExistByWorkspaceThreadChatId(ctx context.Context, workspaceId string, threadChatId string) (bool, error)
	AddLabel(ctx context.Context, thl domain.ThreadChatLabel) (domain.ThreadChatLabel, bool, error)
	GetLabelListByThreadChatId(ctx context.Context, threadChatId string) ([]domain.ThreadChatLabelled, error)
	CreateCustomerThChatMessage(ctx context.Context, threadChatId string, customerId string, msg string) (domain.ThreadChatMessage, error)
	CreateMemberThChatMessage(ctx context.Context, threadChatId string, memberId string, msg string) (domain.ThreadChatMessage, error)
	GetMessageListByThreadChatId(ctx context.Context, threadChatId string) ([]domain.ThreadChatMessage, error)
}
