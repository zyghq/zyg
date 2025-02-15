package services

import (
	"context"

	"github.com/zyghq/zyg/models"

	"github.com/zyghq/zyg/ports"
)

type UserService struct {
	userRepo ports.UserRepositorer
}

func NewUserService(userRepo ports.UserRepositorer) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateWorkOSUser(ctx context.Context, user *models.WorkOSUser) (*models.WorkOSUser, error) {
	user, err := s.userRepo.SaveWorkOSUser(ctx, user)
	if err != nil {
		return &models.WorkOSUser{}, ErrUser
	}
	return user, nil
}
