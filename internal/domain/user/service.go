package user

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
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
		return nil, status.Error(codes.NotFound, ErrUserNotFound)
	}
	if user.Device.IsBlacklisted {
		u.log.DomainUsecase(domainName, "Login").Errorf("user device is blacklisted")
		return nil, status.Error(codes.Forbidden, ErrDeviceIsBlacklisted)
	}
	if !user.Device.Valid(input.Device.FirebaseId, input.Device.DeviceId) {
		u.log.DomainUsecase(domainName, "Login").Errorf("invalid device credentials")
		return nil, status.Error(codes.Forbidden, ErrInvalidDevice)
	}

	matched := u.passwordHasher.Compare(input.Password, user.Password)
	if !matched {
		u.log.DomainUsecase(domainName, "Login").Errorf("invalid password")
		return nil, status.Error(codes.BadRequest, ErrInvalidPassword)
	}

	token, err := u.tokenService.Create(input, tokenExpiredTime)
	if err != nil {
		u.log.DomainUsecase(domainName, "Login").Errorf("Create token: %v", err)
		return nil, status.Error(codes.Internal, ErrCreateTokenFailed)
	}

	return token, nil
}
