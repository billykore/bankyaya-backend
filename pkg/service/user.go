package service

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/pkg/entity"
	pkgerrors "go.bankyaya.org/app/backend/pkg/errors"
	"go.bankyaya.org/app/backend/pkg/interface/repository"
	"go.bankyaya.org/app/backend/pkg/interface/security"
	"go.bankyaya.org/app/backend/pkg/util/codes"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	"go.bankyaya.org/app/backend/pkg/util/status"
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

func NewUser(repo repository.UserRepository, passwordHasher security.PasswordHasher, tokenService security.TokenService) *User {
	return &User{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenService:   tokenService,
	}
}

func (u *User) Login(ctx context.Context, input *entity.User) (*entity.Token, error) {
	user, err := u.repo.GetUserDataByPhoneNumber(ctx, input.PhoneNumber)
	if err != nil {
		u.log.ServiceUsecase(userService, "Login").Errorf("GetUserDataByPhoneNumber: %v", err)
		return nil, status.Error(codes.NotFound, pkgerrors.ErrUserNotFound)
	}
	if user.AuthData.Device.IsBlacklisted() {
		u.log.ServiceUsecase(userService, "Login").Errorf("user is blacklisted")
		return nil, status.Error(codes.Forbidden, pkgerrors.ErrDeviceIsBlacklisted)
	}
	if !user.AuthData.ValidFirebaseId(input.AuthData.FirebaseId) && !user.AuthData.ValidDeviceId(input.AuthData.DeviceId) {
		return nil, pkgerrors.ErrInvalidDevice
	}

	matched := u.passwordHasher.Compare(input.AuthData.Password, user.AuthData.Password)
	if !matched {
		u.log.ServiceUsecase(userService, "Login").Errorf("invalid password")
		return nil, status.Error(codes.BadRequest, pkgerrors.ErrInvalidPassword)
	}

	token, err := u.tokenService.Create(input, tokenExpiredTime)
	if err != nil {
		u.log.ServiceUsecase(userService, "Login").Errorf("Create token: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrCreateTokenFailed)
	}

	return &token, nil
}
