package service

import (
	"time"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/pkg"

	"golang.org/x/crypto/bcrypt"
)

type RefreshToken interface {
	Create(domain.RefreshToken) error
	UpdateRevoke(string) error
	Delete(string) error
	Find(string) (*domain.RefreshToken, error)
}

type AuthService struct {
	repo     RefreshToken
	userRepo UserRepository
}

func NewAuthService(repo RefreshToken, userRepo UserRepository) *AuthService {
	return &AuthService{repo: repo, userRepo: userRepo}
}

func (s *AuthService) SignUp(user *dto.UserReqDTO) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_user := domain.User{
		Username: user.Username,
		Password: string(hashPassword),
	}
	_, err = s.userRepo.Create(_user)
	if err != nil {
		return err
	}

	return nil

}

func (s *AuthService) SignIn(user *dto.UserReqDTO) (*dto.AuthRespDTO, error) {
	_user, err := s.userRepo.FindByUsername(user.Username)
	if err != nil {
		return nil, pkg.ErrRowsEmpty
	}

	if err := bcrypt.CompareHashAndPassword([]byte(_user.Password), []byte(user.Password)); err != nil {
		return nil, pkg.ErrInvalidPassOrUsn
	}

	refToken, err := pkg.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	createRefToken := domain.RefreshToken{
		UserId:    _user.ID,
		TokenHash: refToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3),
	}

	if err := s.repo.Create(createRefToken); err != nil {
		return nil, err
	}

	token, err := pkg.GenerateJWT(_user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthRespDTO{AccessToken: token, RefreshToken: refToken}, nil
}

func (s *AuthService) RefreshToken(token string) (*dto.AuthRespDTO, error) {
	var refToken domain.RefreshToken

	verfToken, err := s.repo.Find(token)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateRevoke(verfToken.TokenHash); err != nil {
		return nil, err
	}

	newToken, err := pkg.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refToken = domain.RefreshToken{
		UserId:    verfToken.UserId,
		TokenHash: newToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 3),
	}

	err = s.repo.Create(refToken)
	if err != nil {
		return nil, err
	}

	newAccToken, err := pkg.GenerateJWT(verfToken.UserId)
	if err != nil {
		return nil, err
	}

	return &dto.AuthRespDTO{AccessToken: newAccToken, RefreshToken: newToken}, nil
}
