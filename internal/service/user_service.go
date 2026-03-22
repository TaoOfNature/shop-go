package service

import (
	"context"
	"errors"

	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/repository"
)

type UserService struct {
	users *repository.UserRepository
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	return s.users.UpdateProfile(ctx, &model.User{
		ID:        userID,
		Nickname:  req.Nickname,
		AvatarURL: req.AvatarURL,
	})
}
