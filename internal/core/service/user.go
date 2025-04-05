package service

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/repository"
	"go.bankyaya.org/app/backend/internal/core/port/security"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

const userService = "User"

const tokenExpiredTime = 15 * time.Minute

// User handles user related process.
type User struct {
	log            *logger.Logger
	repo           repository.UserRepository
	passwordHasher security.PasswordHasher
	tokenService   security.TokenService
}

func NewUser(log *logger.Logger, repo repository.UserRepository, passwordHasher security.PasswordHasher, tokenService security.TokenService) *User {
	return &User{
		log:            log,
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (u *User) Login(ctx context.Context, input *entity.User) (*entity.Token, error) {
	user, err := u.repo.GetUserDataByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		u.log.ServiceUsecase(userService, "Login").Errorf("GetUserDataByPhoneNumber: %v", err)
		return nil, status.Error(codes.NotFound, domain.ErrUserNotFound)
	}
	if user.AuthData.Device.IsBlacklisted() {
		u.log.ServiceUsecase(userService, "Login").Errorf("user device is blacklisted")
		return nil, status.Error(codes.Forbidden, domain.ErrDeviceIsBlacklisted)
	}
	if !user.AuthData.ValidFirebaseId(input.AuthData.FirebaseId) && !user.AuthData.ValidDeviceId(input.AuthData.DeviceId) {
		u.log.ServiceUsecase(userService, "Login").Errorf("invalid device credentials")
		return nil, status.Error(codes.Forbidden, domain.ErrInvalidDevice)
	}

	matched := u.passwordHasher.Compare(input.AuthData.Password, user.AuthData.Password)
	if !matched {
		u.log.ServiceUsecase(userService, "Login").Errorf("invalid password")
		return nil, status.Error(codes.BadRequest, domain.ErrInvalidPassword)
	}

	token, err := u.tokenService.Create(input, tokenExpiredTime)
	if err != nil {
		u.log.ServiceUsecase(userService, "Login").Errorf("Create token: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrCreateTokenFailed)
	}

	return &token, nil
}
