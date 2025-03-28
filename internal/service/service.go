package service

import (
	"context"
	"errors"

	"gravitum/internal/domain/user"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(context.Context, user.User) error
	Update(context.Context, user.User) error
	Delete(context.Context, user.User) error
	GetByID(context.Context, uuid.UUID) (user.User, error)
}

type Service struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

type CreateUserRequest struct {
	Login     string
	FirstName string
	LastName  string
}

type CreateUserResponse struct {
	ID uuid.UUID
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (CreateUserResponse, error) {
	newUser, err := user.New(
		uuid.Must(uuid.NewV7()),
		req.Login,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		return CreateUserResponse{}, err
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return CreateUserResponse{}, err
	}

	return CreateUserResponse{ID: newUser.ID()}, nil
}

type UpdateUserRequest struct {
	ID        string
	Login     string
	FirstName string
	LastName  string
}

func (s *Service) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	u, err := s.GetUserByID(ctx, req.ID)
	if err != nil {
		return err
	}

	updatedUser, err := user.New(
		u.ID(),
		req.Login,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		return err
	}

	return s.userRepo.Update(ctx, updatedUser)
}

func (s *Service) DeleteUserByID(ctx context.Context, id string) error {
	u, err := s.GetUserByID(ctx, id)
	if errors.Is(err, user.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, u)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (user.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return user.User{}, user.ErrID
	}

	return s.userRepo.GetByID(ctx, userID)
}
