package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/adapters/repository"
	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type ThreadChatService struct {
	repo ports.ThreadChatRepositorer
}

func NewThreadChatService(repo ports.ThreadChatRepositorer) *ThreadChatService {
	return &ThreadChatService{
		repo: repo,
	}
}

// creates a new thread chat for the customer and message.
// a thread can only be created for a customer.
// thread cannot exists without a valid customer.
func (s *ThreadChatService) CreateThreadWithMessage(ctx context.Context, th models.ThreadChat, msg string,
) (models.ThreadChat, models.ThreadChatMessage, error) {
	thread, message, err := s.repo.InsertThreadChat(ctx, th, msg)

	if err != nil {
		return thread, message, ErrThreadChat
	}

	return thread, message, nil
}

func (s *ThreadChatService) UpdateThread(ctx context.Context, th models.ThreadChat, fields []string) (models.ThreadChat, error) {
	thread, err := s.repo.ModifyThreadChatById(ctx, th, fields)

	if err != nil {
		return thread, ErrThreadChat
	}

	return thread, nil
}

// returns a thread chat by workspace and thread chat.
// a thread chat is unique in a workspace.
func (s *ThreadChatService) GetThread(ctx context.Context, workspaceId string, threadChatId string,
) (models.ThreadChat, error) {
	thread, err := s.repo.LookupByWorkspaceThreadChatId(ctx, workspaceId, threadChatId)

	if errors.Is(err, repository.ErrEmpty) {
		return thread, ErrThreadChatNotFound
	}

	if err != nil {
		return thread, ErrThreadChat
	}

	return thread, nil
}

// returns a list of thread chat for the customer in the workspace.
func (s *ThreadChatService) ListCustomerThreads(ctx context.Context, workspaceId string, customerId string,
) ([]models.ThreadChatWithMessage, error) {
	threads, err := s.repo.FetchByWorkspaceCustomerId(ctx, workspaceId, customerId)

	if err != nil {
		return threads, ErrThreadChat
	}

	return threads, nil
}

// assigns a member to a thread chat.
func (s *ThreadChatService) AssignMemberToThread(ctx context.Context, threadChatId string, assigneeId string,
) (models.ThreadChat, error) {
	thread, err := s.repo.UpdateAssignee(ctx, threadChatId, assigneeId)
	if err != nil {
		return thread, ErrThreadChatAssign
	}
	return thread, nil
}

// marks a thread chat as replied or unreplied.
func (s *ThreadChatService) SetThreadReplyStatus(ctx context.Context, threadChatId string, replied bool,
) (models.ThreadChat, error) {
	thread, err := s.repo.UpdateRepliedStatus(ctx, threadChatId, replied)
	if err != nil {
		return thread, ErrThreadChatReplied
	}
	return thread, nil
}

// returns a list of thread chat in the workspace.
func (s *ThreadChatService) ListWorkspaceThreads(ctx context.Context, workspaceId string) ([]models.ThreadChatWithMessage, error) {
	threads, err := s.repo.RetrieveByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return threads, ErrThreadChat
	}
	return threads, nil
}

// returns a list of thread chat assigned to the member in the workspace.
func (s *ThreadChatService) ListMemberAssignedThreads(ctx context.Context, workspaceId string, memberId string,
) ([]models.ThreadChatWithMessage, error) {
	threads, err := s.repo.FetchAssignedThreadsByMember(ctx, workspaceId, memberId)
	if err != nil {
		return threads, ErrThreadChat
	}
	return threads, nil
}

// returns a list of thread chat unassigned in the workspace.
func (s *ThreadChatService) ListUnassignedThreads(ctx context.Context, workspaceId string,
) ([]models.ThreadChatWithMessage, error) {
	threads, err := s.repo.RetrieveUnassignedThreads(ctx, workspaceId)
	if err != nil {
		return threads, ErrThreadChat
	}
	return threads, nil
}

// returns a list of thread chat labelled in the workspace.
func (s *ThreadChatService) ListLabelledThreads(ctx context.Context, workspaceId string, labelId string,
) ([]models.ThreadChatWithMessage, error) {
	threads, err := s.repo.FetchThreadsByLabel(ctx, workspaceId, labelId)
	if err != nil {
		return threads, ErrThreadChat
	}
	return threads, nil
}

// checks if a thread chat exists in the workspace.
func (s *ThreadChatService) ThreadExistsInWorkspace(ctx context.Context, workspaceId string, threadChatId string,
) (bool, error) {
	exist, err := s.repo.CheckExistenceByWorkspaceThreadChatId(ctx, workspaceId, threadChatId)
	if err != nil {
		return exist, ErrThreadChat
	}
	return exist, nil
}

// adds a label to a thread chat.
// label must exist in the workspace.
func (s *ThreadChatService) AttachLabelToThread(ctx context.Context, thl models.ThreadChatLabel,
) (models.ThreadChatLabel, bool, error) {
	label, created, err := s.repo.AttachLabelToThread(ctx, thl)
	if err != nil {
		return label, created, ErrThChatLabel
	}

	return label, created, nil
}

// returns a list of labels added to the thread chat.
func (s *ThreadChatService) ListThreadLabels(ctx context.Context, threadChatId string) ([]models.ThreadChatLabelled, error) {
	labels, err := s.repo.RetrieveLabelsByThreadChatId(ctx, threadChatId)
	if err != nil {
		return labels, ErrThChatLabel
	}
	return labels, nil
}

// creates a new message for the customer in the thread chat.
// a customer thread chat must exist before creating a message.
func (s *ThreadChatService) AddCustomerMessageToThread(ctx context.Context, th models.ThreadChat, c *models.Customer, msg string,
) (models.ThreadChatMessage, error) {
	message, err := s.repo.InsertCustomerMessage(ctx, th.ThreadChatId, c.CustomerId, msg)
	if err != nil {
		return models.ThreadChatMessage{}, ErrThChatMessage
	}

	return message, nil
}

// creates a new message for the member in the thread chat.
// a member thread chat must exist before creating a message.
func (s *ThreadChatService) AddMemberMessageToThread(ctx context.Context, th models.ThreadChat, m *models.Member, msg string,
) (models.ThreadChatMessage, error) {
	message, err := s.repo.InsertMemberMessage(ctx, th.ThreadChatId, m.MemberId, msg)
	if err != nil {
		return models.ThreadChatMessage{}, ErrThChatMessage
	}

	return message, nil
}

// returns a list of messages for the thread chat item.
func (s *ThreadChatService) ListThreadMessages(ctx context.Context, threadChatId string,
) ([]models.ThreadChatMessage, error) {
	messages, err := s.repo.FetchMessagesByThreadChatId(ctx, threadChatId)
	if err != nil {
		return messages, ErrThChatMessage
	}

	return messages, nil
}

// generates metrics for the member in the workspace.
func (s *ThreadChatService) GenerateMemberThreadMetrics(ctx context.Context, workspaceId string, memberId string,
) (models.ThreadMemberMetrics, error) {
	statusMetrics, err := s.repo.ComputeStatusMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	assignmentMetrics, err := s.repo.CalculateAssigneeMetricsByMember(ctx, workspaceId, memberId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	labelMetrics, err := s.repo.ComputeLabelMetricsByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return models.ThreadMemberMetrics{}, ErrThreadChatMetrics
	}

	metrics := models.ThreadMemberMetrics{
		ThreadMetrics:         statusMetrics,
		ThreadAssigneeMetrics: assignmentMetrics,
		ThreadLabelMetrics:    labelMetrics,
	}

	return metrics, nil
}
