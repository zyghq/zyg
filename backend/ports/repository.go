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

	InsertSystemMember(
		ctx context.Context, member models.Member) (models.Member, error)
	LookupSystemMemberByOldest(
		ctx context.Context, workspaceId string) (models.Member, error)
	SavePostmarkSetting(
		ctx context.Context, setting models.PostmarkServerSetting) (models.PostmarkServerSetting, error)
	FetchPostmarkSettingById(
		ctx context.Context, workspaceId string) (models.PostmarkServerSetting, error)
	ModifyPostmarkSettingById(
		ctx context.Context, setting models.PostmarkServerSetting, fields []string,
	) (models.PostmarkServerSetting, error)
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

	FetchCustomersByWorkspaceId(
		ctx context.Context, workspaceId string, role *string) ([]models.Customer, error)

	ModifyCustomerById(
		ctx context.Context, customer models.Customer) (models.Customer, error)

	InsertEvent(ctx context.Context, event models.Event) (models.Event, error)
	FetchEventsByCustomerId(ctx context.Context, customerId string) ([]models.Event, error)
}

type ThreadRepositorer interface {
	// SaveThreadActivity persists thread and activity in transaction,
	// upserts the thread, and inserts new activity.
	SaveThreadActivity(
		ctx context.Context, thread *models.Thread, activity *models.Activity) (*models.Thread, *models.Activity, error)

	// SavePostmarkThreadActivity persists thread, activity, and postmark message log,
	// upserts the thread, inserts new activity and postmark message log.
	SavePostmarkThreadActivity(
		ctx context.Context, thread *models.Thread, activity *models.Activity,
		postmarkMessageLog *models.PostmarkMessageLog) (*models.Thread, *models.Activity, error)

	CheckPostmarkInboundExists(ctx context.Context, pmMessageId string) (bool, error)

	InsertMessageAttachment(
		ctx context.Context, message models.ActivityAttachment) (models.ActivityAttachment, error)

	FetchMessageAttachmentById(
		ctx context.Context, messageId, attachmentId string) (models.ActivityAttachment, error)

	FindThreadByPostmarkReplyMessageId(
		ctx context.Context, workspaceId string, inReplyMessageId string) (models.Thread, error)

	GetRecentMailMessageIdByThreadId(ctx context.Context, threadId string) (string, error)

	LookupByWorkspaceThreadId(
		ctx context.Context, workspaceId string, threadId string, channel *string) (models.Thread, error)
	ModifyThreadById(
		ctx context.Context, thread models.Thread, fields []string) (models.Thread, error)

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
		ctx context.Context, threadId string) ([]models.Activity, error)

	FetchMessagesWithAttachmentsByThreadId(
		ctx context.Context, threadId string) ([]models.ActivityWithAttachments, error)

	ComputeStatusMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) (models.ThreadMetrics, error)
	ComputeAssigneeMetricsByMember(
		ctx context.Context, workspaceId string, memberId string) (models.ThreadAssigneeMetrics, error)
	ComputeLabelMetricsByWorkspaceId(
		ctx context.Context, workspaceId string) ([]models.ThreadLabelMetric, error)
	DeleteThreadLabelById(
		ctx context.Context, threadId string, labelId string) error
}
