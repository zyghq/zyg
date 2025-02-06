package ports

import (
	"context"
	"github.com/zyghq/zyg/models"
)

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
	SavePostmarkMailServerSetting(
		ctx context.Context, setting models.PostmarkMailServerSetting) (models.PostmarkMailServerSetting, error)
	FetchPostmarkMailServerSettingById(
		ctx context.Context, workspaceId string) (models.PostmarkMailServerSetting, error)
	ModifyPostmarkMailServerSettingById(
		ctx context.Context, setting models.PostmarkMailServerSetting, fields []string,
	) (models.PostmarkMailServerSetting, error)
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
	LookupWorkspaceCustomerByEmail(
		ctx context.Context, workspaceId string, email string, role *string) (models.Customer, error)
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
	InsertClaimedMail(
		ctx context.Context, claimed models.ClaimedMail) (models.ClaimedMail, error)
	DeleteCustomerClaimedMail(
		ctx context.Context, workspaceId string, customerId string, email string) error
	LookupClaimedMailByToken(
		ctx context.Context, token string) (models.ClaimedMail, error)
	LookupLatestClaimedMail(
		ctx context.Context, workspaceId string, customerId string,
	) (models.ClaimedMail, error)
	InsertEvent(ctx context.Context, event models.Event) (models.Event, error)
	FetchEventsByCustomerId(ctx context.Context, customerId string) ([]models.Event, error)
}

type ThreadRepositorer interface {
	InsertInboundThreadMessage(
		ctx context.Context, thread *models.Thread, message *models.Message) (*models.Thread, *models.Message, error)
	AppendInboundThreadMessage(
		ctx context.Context, thread *models.Thread, message *models.Message) (*models.Thread, *models.Message, error)

	InsertPostmarkInboundThreadMessage(
		ctx context.Context, thread *models.Thread, message *models.Message,
		postmarkMessageLog *models.PostmarkMessageLog) (*models.Thread, *models.Message, error)
	AppendPostmarkInboundThreadMessage(
		ctx context.Context, thread *models.Thread, message *models.Message,
		postmarkMessageLog *models.PostmarkMessageLog) (*models.Message, error)

	AppendOutboundThreadMessage(
		ctx context.Context, thread *models.Thread, message *models.Message) (*models.Thread, *models.Message, error)

	CheckPostmarkInboundMessageExists(ctx context.Context, messageId string) (bool, error)

	InsertMessageAttachment(
		ctx context.Context, message models.MessageAttachment) (models.MessageAttachment, error)

	FetchMessageAttachmentById(
		ctx context.Context, messageId, attachmentId string) (models.MessageAttachment, error)

	FindThreadByPostmarkReplyMessageId(
		ctx context.Context, workspaceId string, inReplyMessageId string) (models.Thread, error)

	GetRecentMailMessageIdByThreadId(ctx context.Context, threadId string) (string, error)

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

	FetchMessagesByThreadId(
		ctx context.Context, threadId string) ([]models.Message, error)

	FetchMessagesWithAttachmentsByThreadId(
		ctx context.Context, threadId string) ([]models.MessageWithAttachments, error)

	ComputeStatusMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	ComputeAssigneeMetricsByMember(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	ComputeLabelMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
	DeleteThreadLabelById(
		ctx context.Context, threadId string, labelId string) error
}
