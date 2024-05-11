package services

import (
	"context"
	"errors"

	"github.com/zyghq/zyg/internal/adapters/repository"
	"github.com/zyghq/zyg/internal/domain"
	"github.com/zyghq/zyg/internal/ports"
)

type ThreadChatService struct {
	repo ports.ThreadChatRepositorer
}

func NewThreadChatService(repo ports.ThreadChatRepositorer) *ThreadChatService {
	return &ThreadChatService{
		repo: repo,
	}
}

func (s *ThreadChatService) CreateCustomerThread(ctx context.Context, th domain.ThreadChat, msg string,
) (domain.ThreadChat, domain.ThreadChatMessage, error) {
	thread, message, err := s.repo.CreateThreadChat(ctx, th, msg)

	if errors.Is(err, repository.ErrTxQuery) {
		return thread, message, ErrThreadChat
	}

	if errors.Is(err, repository.ErrQuery) {
		return thread, message, ErrThreadChat
	}

	if errors.Is(err, repository.ErrEmpty) {
		return thread, message, ErrThreadChat
	}

	if err != nil {
		return thread, message, err
	}

	return thread, message, nil
}

func (s *ThreadChatService) GetWorkspaceThread(ctx context.Context, workspaceId string, threadChatId string,
) (domain.ThreadChat, error) {
	thread, err := s.repo.GetByWorkspaceThreadChatId(ctx, workspaceId, threadChatId)

	if errors.Is(err, repository.ErrEmpty) {
		return thread, ErrThreadChatNotFound
	}

	if err != nil {
		return thread, ErrThreadChat
	}

	return thread, nil
}

func (s *ThreadChatService) GetWorkspaceCustomerList(ctx context.Context, workspaceId string, customerId string,
) ([]domain.ThreadChatWithMessage, error) {
	threads, err := s.repo.GetListByWorkspaceCustomerId(ctx, workspaceId, customerId)

	if errors.Is(err, repository.ErrQuery) {
		return threads, ErrThreadChat
	}

	if err != nil {
		return threads, err
	}
	return threads, nil
}

func (s *ThreadChatService) AssignMember(ctx context.Context, threadChatId string, assigneeId string,
) (domain.ThreadChat, error) {
	thread, err := s.repo.SetAssignee(ctx, threadChatId, assigneeId)
	if err != nil {
		return thread, err
	}
	return thread, nil
}

func (s *ThreadChatService) MarkReplied(ctx context.Context, threadChatId string, replied bool,
) (domain.ThreadChat, error) {
	thread, err := s.repo.SetReplied(ctx, threadChatId, replied)
	if err != nil {
		return thread, err
	}
	return thread, nil
}

func (s *ThreadChatService) GetWorkspaceList(ctx context.Context, workspaceId string,
) ([]domain.ThreadChatWithMessage, error) {
	threads, err := s.repo.GetListByWorkspaceId(ctx, workspaceId)
	if err != nil {
		return threads, err
	}
	return threads, nil
}

func (s *ThreadChatService) ExistInWorkspace(ctx context.Context, workspaceId string, threadChatId string,
) (bool, error) {
	exist, err := s.repo.IsExistByWorkspaceThreadChatId(ctx, workspaceId, threadChatId)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

func (s *ThreadChatService) AddLabel(ctx context.Context, thl domain.ThreadChatLabel,
) (domain.ThreadChatLabel, bool, error) {
	label, created, err := s.repo.AddLabel(ctx, thl)

	if errors.Is(err, repository.ErrQuery) {
		return label, false, ErrThChatLabel
	}

	if errors.Is(err, repository.ErrEmpty) {
		return label, false, ErrThChatLabelNotFound
	}

	if err != nil {
		return label, created, err
	}

	return label, created, nil
}

func (s *ThreadChatService) GetLabelList(ctx context.Context, threadChatId string) ([]domain.ThreadChatLabelled, error) {
	labels, err := s.repo.GetLabelListByThreadChatId(ctx, threadChatId)
	if err != nil {
		return labels, err
	}
	return labels, nil
}

func (s *ThreadChatService) CreateCustomerMessage(ctx context.Context, th domain.ThreadChat, c *domain.Customer, msg string,
) (domain.ThreadChatMessage, error) {
	message, err := s.repo.CreateCustomerThChatMessage(ctx, th.ThreadChatId, c.CustomerId, msg)

	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return domain.ThreadChatMessage{}, ErrThreadChat
	}

	if err != nil {
		return domain.ThreadChatMessage{}, err
	}

	return message, nil
}

func (s *ThreadChatService) CreateMemberMessage(ctx context.Context, th domain.ThreadChat, m *domain.Member, msg string,
) (domain.ThreadChatMessage, error) {
	message, err := s.repo.CreateMemberThChatMessage(ctx, th.ThreadChatId, m.MemberId, msg)

	if errors.Is(err, repository.ErrQuery) || errors.Is(err, repository.ErrEmpty) {
		return domain.ThreadChatMessage{}, ErrThreadChat
	}

	if err != nil {
		return domain.ThreadChatMessage{}, err
	}

	return message, nil
}

func (s *ThreadChatService) GetMessageList(ctx context.Context, threadChatId string,
) ([]domain.ThreadChatMessage, error) {
	messages, err := s.repo.GetMessageListByThreadChatId(ctx, threadChatId)

	if errors.Is(err, repository.ErrQuery) {
		return messages, ErrThChatMessage
	}

	if err != nil {
		return messages, err
	}

	return messages, nil
}
