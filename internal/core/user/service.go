package user

import (
	"context"
	"time"
)

const tokenExpiredTime = 15 * time.Minute

// Service handles user related process.
type Service struct {
	repo           Repository
	passwordHasher PasswordHasher
	tokenService   TokenService
}

func NewService(repo Repository, passwordHasher PasswordHasher, tokenService TokenService) *Service {
	return &Service{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (s *Service) Login(ctx context.Context, input *User) (*Token, error) {
	user, err := s.repo.GetUserDataByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if user.AuthData.Device.IsBlacklisted() {
		return nil, ErrDeviceIsBlacklisted
	}
	if user.AuthData.ValidFirebaseId(input.AuthData.FirebaseId) && user.AuthData.ValidDeviceId(input.AuthData.DeviceId) {
		return nil, ErrInvalidDevice
	}

	matched := s.passwordHasher.Compare(input.AuthData.Password, user.AuthData.Password)
	if !matched {
		return nil, ErrInvalidPassword
	}

	token, err := s.tokenService.Create(input, tokenExpiredTime)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
