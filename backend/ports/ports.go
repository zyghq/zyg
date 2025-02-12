package ports

import (
	"context"
	"time"

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

//
//type CustomerAuthServicer interface {
//	AuthenticateWorkspaceCustomer(
//		ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error)
//	GetWidgetLinkedSecretKey(
//		ctx context.Context, widgetId string) (models.WorkspaceSecret, error)
//}

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
		ctx context.Context, workspaceId string, externalId string, name string) (models.Customer, bool, error)
	CreateCustomerWithEmail(
		ctx context.Context, workspaceId string, email string, isEmailVerified bool, name string,
	) (models.Customer, bool, error)
	CreateCustomerWithPhone(
		ctx context.Context, workspaceId string, phone string, name string,
	) (models.Customer, bool, error)
	CreateWidget(
		ctx context.Context, workspaceId string, name string, configuration map[string]interface{},
	) (models.Widget, error)
	ListWidgets(ctx context.Context, workspaceId string) ([]models.Widget, error)
	GenerateWorkspaceSecret(
		ctx context.Context, workspaceId string, length int) (models.WorkspaceSecret, error)
	GetSecretKey(
		ctx context.Context, workspaceId string) (models.WorkspaceSecret, error)
	GetCustomer(
		ctx context.Context, workspaceId string, customerId string, role *string) (models.Customer, error)

	PostmarkCreateMailServer(
		ctx context.Context, workspaceId, email, domain string) (models.PostmarkServerSetting, error)
	GetPostmarkMailServerSetting(
		ctx context.Context, workspaceId string) (models.PostmarkServerSetting, error)
	PostmarkMailServerAddDomain(
		ctx context.Context, setting models.PostmarkServerSetting, domain string,
	) (models.PostmarkServerSetting, bool, error)
	PostmarkMailServerVerifyDomain(
		ctx context.Context, setting models.PostmarkServerSetting) (models.PostmarkServerSetting, error)
	PostmarkMailServerUpdate(
		ctx context.Context, setting models.PostmarkServerSetting, fields []string,
	) (models.PostmarkServerSetting, error)
}

type CustomerServicer interface {
	GenerateCustomerJwt(
		customer models.Customer, sk string) (string, error)
	VerifyExternalId(sk string, hash string, externalId string) bool
	VerifyEmail(sk string, hash string, email string) bool
	VerifyPhone(sk string, hash string, phone string) bool
	UpdateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error)

	GenerateMailVerificationToken(
		sk string, workspaceId string, customerId string, email string,
		expiresAt time.Time, redirectUrl string,
	) (string, error)

	AddEvent(ctx context.Context, event models.Event) (models.Event, error)

	ListEvents(ctx context.Context, customerId string) ([]models.Event, error)
}

type ThreadServicer interface {
	GetRecentThreadMailMessageId(ctx context.Context, threadId string) (string, error)
	IsPostmarkInboundProcessed(ctx context.Context, pmMessageId string) (bool, error)

	ProcessPostmarkInbound(
		ctx context.Context, workspaceId string,
		customer *models.Customer, createdBy *models.Member, inboundMessage *models.PostmarkInboundMessage,
	) (models.Thread, models.Activity, error)
	SendThreadMailReply(
		ctx context.Context,
		workspace *models.Workspace, setting *models.PostmarkServerSetting, thread *models.Thread,
		member *models.Member, customer *models.Customer,
		textBody, htmlBody string,
	) (models.Thread, models.Activity, error)

	GetPostmarkInReplyThread(
		ctx context.Context, workspaceId, mailMessageId string) (*models.Thread, error)

	GetWorkspaceThread(
		ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error)
	UpdateThread(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)

	ListWorkspaceThreads(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListMemberThreads(
		ctx context.Context, memberId string) ([]models.Thread, error)
	ListUnassignedThreads(
		ctx context.Context, workspaceId string) ([]models.Thread, error)
	ListLabelledThreads(
		ctx context.Context, labelId string) ([]models.Thread, error)

	ThreadExistsInWorkspace(
		ctx context.Context, workspaceId string, threadId string) (bool, error)

	SetLabel(
		ctx context.Context, threadId string, labelId string, addedBy string) (models.ThreadLabel, bool, error)
	ListThreadLabels(
		ctx context.Context, threadId string) ([]models.ThreadLabel, error)
	RemoveThreadLabel(
		ctx context.Context, threadId string, labelId string) error

	ListThreadMessageActivities(ctx context.Context, threadId string) ([]models.Activity, error)
	ListThreadMessagesWithAttachments(
		ctx context.Context, threadId string) ([]models.ActivityWithAttachments, error)

	GetMessageAttachment(
		ctx context.Context, messageId, attachmentId string) (models.ActivityAttachment, error)

	GenerateMemberThreadMetrics(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadMemberMetrics, error)

	LogPostmarkInboundRequest(
		ctx context.Context, workspaceId, messageId string, payload map[string]interface{}) error
}
