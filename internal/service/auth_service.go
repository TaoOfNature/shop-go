package service

import (
	"context"
	"errors"

	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/pkg/auth"
	"github.com/dawnstack/shop-go/internal/pkg/password"
	"github.com/dawnstack/shop-go/internal/pkg/snowflake"
	"github.com/dawnstack/shop-go/internal/repository"
)

type AuthService struct {
	users  *repository.UserRepository
	idGen  *snowflake.Generator
	tokens *auth.Manager
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthTokens struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

func NewAuthService(users *repository.UserRepository, idGen *snowflake.Generator, tokens *auth.Manager) *AuthService {
	return &AuthService{users: users, idGen: idGen, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthTokens, error) {
	existing, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           s.idGen.NextID(),
		Email:        req.Email,
		PasswordHash: hashed,
		Nickname:     req.Nickname,
	}
	if user.Nickname == "" {
		user.Nickname = "user_" + req.Email
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.tokens.GeneratePair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken, User: user}, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthTokens, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := password.Compare(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := s.tokens.GeneratePair(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	return &AuthTokens{AccessToken: accessToken, RefreshToken: refreshToken, User: user}, nil
}

func (s *AuthService) Refresh(_ context.Context, req RefreshRequest) (*AuthTokens, error) {
	claims, err := s.tokens.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, refreshToken, err := s.tokens.GeneratePair(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}
	return &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &model.User{
			ID:    claims.UserID,
			Email: claims.Email,
		},
	}, nil
}

func (s *AuthService) TokenManager() *auth.Manager {
	return s.tokens
}
