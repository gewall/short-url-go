package service

import (
	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(domain.User) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByID(id uuid.UUID) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Delete(id uuid.UUID) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user dto.UserReqDTO) (*dto.UserRespDTO, error) {
	u := domain.User{
		Username: user.Username,
		Password: user.Password,
	}

	created, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}
	return &dto.UserRespDTO{
		ID:        created.ID,
		Username:  created.Username,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s *UserService) FindUserByID(id uuid.UUID) (*dto.UserRespDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, pkg.ErrRowsEmpty
	}
	return &dto.UserRespDTO{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) FindAllUsers() ([]*dto.UserRespDTO, error) {
	var result []*dto.UserRespDTO
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		result = append(result, &dto.UserRespDTO{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		})
	}

	return result, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
