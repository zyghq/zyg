package services

import (
	"context"
	"fmt"

	"github.com/zyghq/zyg/models"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateWorkOSUser(ctx context.Context, user *models.WorkOSUser) (*models.WorkOSUser, error) {
	fmt.Println("SERVICES: Creating WorkOS user", user)
	return user, nil
}
