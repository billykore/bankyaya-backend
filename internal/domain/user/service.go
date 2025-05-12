package user

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/pkgerror"
)

const (
	domainName       = "user"
	tokenExpiredTime = 15 * time.Minute
)

// Service handles user-related process.
type Service struct {
	log            *logger.Logger
	repo           Repository
	passwordHasher PasswordHasher
	tokenService   TokenService
}

func NewService(log *logger.Logger, repo Repository, passwordHasher PasswordHasher, tokenService TokenService) *Service {
	return &Service{
		log:            log,
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (u *Service) Login(ctx context.Context, input *User) (*Token, error) {
	user, err := u.repo.GetUserByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		u.log.DomainUsecase(domainName, "Login").Errorf("GetUserByPhoneNumber: %v", err)
		return nil, pkgerror.New(codes.NotFound, ErrUserNotFound).
			SetMsg("User not found. Please register your account first.")
	}
	if user.Device.IsBlacklisted {
		u.log.DomainUsecase(domainName, "Login").Error(ErrDeviceIsBlacklisted)
		return nil, pkgerror.New(codes.Forbidden, ErrDeviceIsBlacklisted).
			SetMsg("Device is blacklisted. Please contact support.")
	}
	if !user.Device.Valid(input.Device.FirebaseID, input.Device.DeviceID) {
		u.log.DomainUsecase(domainName, "Login").Error(ErrInvalidDevice)
		return nil, pkgerror.New(codes.Forbidden, ErrInvalidDevice).
			SetMsg("Device is not registered. Please register your device first.")
	}

	matched := u.passwordHasher.Compare(input.Password, user.Password)
	if !matched {
		u.log.DomainUsecase(domainName, "Login").Error(ErrInvalidPassword)
		return nil, pkgerror.New(codes.BadRequest, ErrInvalidPassword).
			SetMsg("Password is incorrect. Please try again.")
	}

	token, err := u.tokenService.Create(input, tokenExpiredTime)
	if err != nil {
		u.log.DomainUsecase(domainName, "Login").Errorf("Create token: %v", err)
		return nil, pkgerror.New(codes.Internal, ErrCreateTokenFailed).
			SetMsg("Login failed. Please try again later.")
	}

	return token, nil
}
