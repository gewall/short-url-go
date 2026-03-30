package service_test

import (
	"errors"
	"testing"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
}
type mockRefreshToken struct {
	mock.Mock
}

func (m *mockUserRepository) Create(user domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), nil
}

func (m *mockUserRepository) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindAll() ([]*domain.User, error) {
	args := m.Called()
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *mockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockRefreshToken) Create(token domain.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockRefreshToken) UpdateRevoke(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockRefreshToken) Delete(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockRefreshToken) Find(token string) (*domain.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func TestSignUp_Success(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, mockUserRepo)

	user := dto.UserReqDTO{
		Username: "testuser",
		Password: "password",
	}

	mockUserRepo.
		On("Create", mock.MatchedBy(func(u domain.User) bool {
			return u.Username == user.Username &&
				u.Password != user.Password // karena sudah di-hash
		})).
		Return(&domain.User{
			ID:       uuid.New(),
			Username: user.Username,
		}, nil)

	err := authService.SignUp(&user)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
}

func TestSignUp_Failure(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, mockUserRepo)

	user := dto.UserReqDTO{
		Username: "testuser",
		Password: "password",
	}

	mockUserRepo.
		On("Create", mock.MatchedBy(func(u domain.User) bool {
			return u.Username == user.Username &&
				u.Password != user.Password // karena sudah di-hash
		})).
		Return(nil, errors.New("database error"))

	err := authService.SignUp(&user)

	if err == nil {
		t.Error("SignUp should have failed")
	}

	mockUserRepo.AssertCalled(t, "Create", mock.MatchedBy(func(u domain.User) bool {
		return u.Username == user.Username &&
			u.Password != user.Password // karena password sudah di-hash
	}))
}

func TestSignIn_Success(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, mockUserRepo)

	user := dto.UserReqDTO{
		Username: "testuser",
		Password: "password",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	mockUserRepo.
		On("FindByUsername", user.Username).
		Return(&domain.User{
			ID:       uuid.New(),
			Username: user.Username,
			Password: string(hashedPassword),
		}, nil)

	mockRefresh.
		On("Create", mock.MatchedBy(func(r domain.RefreshToken) bool {
			return r.UserId != uuid.Nil
		})).
		Return(nil)

	token, err := authService.SignIn(&user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)

	mockUserRepo.AssertCalled(t, "FindByUsername", user.Username)
	mockRefresh.AssertCalled(t, "Create", mock.MatchedBy(func(r domain.RefreshToken) bool {
		return r.UserId != uuid.Nil
	}))
}

func TestSignIn_Failure(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, mockUserRepo)

	user := dto.UserReqDTO{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockUserRepo.
		On("FindByUsername", user.Username).
		Return(&domain.User{
			ID:       uuid.New(),
			Username: user.Username,
			Password: string(hashedPassword),
		}, nil)

	mockRefresh.
		On("Create", mock.MatchedBy(func(r domain.RefreshToken) bool {
			return r.UserId != uuid.Nil
		})).
		Return(nil)

	token, err := authService.SignIn(&user)

	assert.Error(t, err)
	assert.Nil(t, token)

	mockUserRepo.AssertCalled(t, "FindByUsername", user.Username)
	mockRefresh.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSignIn_InvalidPassword(t *testing.T) {
	mockUserRepo := new(mockUserRepository)
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, mockUserRepo)

	user := dto.UserReqDTO{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockUserRepo.
		On("FindByUsername", user.Username).
		Return(&domain.User{
			ID:       uuid.New(),
			Username: user.Username,
			Password: string(hashedPassword),
		}, nil)

	mockRefresh.
		On("Create", mock.MatchedBy(func(r domain.RefreshToken) bool {
			return r.UserId != uuid.Nil
		})).
		Return(nil)

	token, err := authService.SignIn(&user)

	assert.Error(t, err)
	assert.Nil(t, token)

	mockUserRepo.AssertCalled(t, "FindByUsername", user.Username)
	mockRefresh.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRefreshToken_Success(t *testing.T) {
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, nil)

	mockRefresh.On("Find", "123").Return(&domain.RefreshToken{
		UserId:    uuid.New(),
		TokenHash: "123",
	}, nil)

	mockRefresh.On("UpdateRevoke", "123").Return(nil)

	mockRefresh.
		On("Create", mock.MatchedBy(func(r domain.RefreshToken) bool {
			return r.UserId != uuid.Nil
		})).
		Return(nil)

	token, err := authService.RefreshToken("123")

	assert.NoError(t, err)
	assert.NotNil(t, token)

	mockRefresh.AssertCalled(t, "Find", "123")
	mockRefresh.AssertCalled(t, "UpdateRevoke", "123")
	mockRefresh.AssertCalled(t, "Create", mock.Anything)

}

func TestRefreshToken_Failure(t *testing.T) {
	mockRefresh := new(mockRefreshToken)
	authService := service.NewAuthService(mockRefresh, nil)

	mockRefresh.On("Find", "123").Return(&domain.RefreshToken{}, errors.New("refresh token not found"))

	token, err := authService.RefreshToken("123")

	assert.Error(t, err)
	assert.Nil(t, token)

	mockRefresh.AssertCalled(t, "Find", "123")
	mockRefresh.AssertNotCalled(t, "UpdateRevoke", "123")
	mockRefresh.AssertNotCalled(t, "Create", mock.Anything)
}
