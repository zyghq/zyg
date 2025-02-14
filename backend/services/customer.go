package services

import (
	"context"

	"github.com/zyghq/zyg/models"
	"github.com/zyghq/zyg/ports"
)

type CustomerService struct {
	repo ports.CustomerRepositorer
}

func NewCustomerService(
	repo ports.CustomerRepositorer) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) AddEvent(
	ctx context.Context, event models.Event) (models.Event, error) {
	event, err := s.repo.InsertEvent(ctx, event)
	if err != nil {
		return models.Event{}, ErrCustomerEvent
	}
	return event, nil
}

func (s *CustomerService) ListEvents(
	ctx context.Context, customerId string) ([]models.Event, error) {
	events, err := s.repo.FetchEventsByCustomerId(ctx, customerId)
	if err != nil {
		return []models.Event{}, ErrCustomerEvent
	}
	return events, nil
}
