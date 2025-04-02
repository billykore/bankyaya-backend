package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/repository/mock"
	securitymock "go.bankyaya.org/app/backend/internal/core/port/security/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

func TestSuccessLogin(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338442777").
		Return(&entity.User{
			AuthData: entity.AuthData{
				Password:   "password",
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.On("Compare", "password", "password").
		Return(true)

	tokenSvcMock.On("Create", mock.Anything, 15*time.Minute).
		Return(entity.Token{
			AccessToken: "example-token-123",
			ExpiredTime: time.Now().Add(15 * time.Minute).Unix(),
		}, nil)

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338442777",
		AuthData: entity.AuthData{
			Password:   "password",
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, token, &entity.Token{
		AccessToken: "example-token-123",
		ExpiredTime: time.Now().Add(15 * time.Minute).Unix(),
	})

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}

func TestLoginFailed_UserNotFound(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338000000").
		Return(nil, errors.New("user not found"))

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338000000",
		AuthData: entity.AuthData{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.NotFound, domain.ErrUserNotFound))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_UserBlacklisted(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338000001").
		Return(&entity.User{
			AuthData: entity.AuthData{
				Device: entity.Device{
					Blacklist: entity.BlacklistDevice{
						Status: "active",
					},
				},
			},
		}, nil)

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338000001",
		AuthData: entity.AuthData{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Forbidden, domain.ErrDeviceIsBlacklisted))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_InvalidDevice(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338000002").
		Return(&entity.User{
			AuthData: entity.AuthData{
				Device: entity.Device{
					FirebaseId: "zxc",
					DeviceId:   "asd",
				},
			},
		}, nil)

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338000002",
		AuthData: entity.AuthData{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Forbidden, domain.ErrInvalidDevice))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_InvalidPassword(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338000003").
		Return(&entity.User{
			AuthData: entity.AuthData{
				Password:   "password",
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.On("Compare", "invalid-password", "password").
		Return(false)

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338000003",
		AuthData: entity.AuthData{
			Password:   "invalid-password",
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.BadRequest, domain.ErrInvalidPassword))

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}

func TestLoginFailed_CreateTokenFailed(t *testing.T) {
	var (
		repoMock     = repomock.NewUserRepository(t)
		hasherMock   = securitymock.NewPasswordHasher(t)
		tokenSvcMock = securitymock.NewTokenService(t)
		svc          = NewUser(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.On("GetUserDataByPhoneNumber", mock.Anything, "081338442777").
		Return(&entity.User{
			AuthData: entity.AuthData{
				Password:   "password",
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.On("Compare", "password", "password").
		Return(true)

	tokenSvcMock.On("Create", mock.Anything, 15*time.Minute).
		Return(entity.Token{}, errors.New("failed to create token"))

	token, err := svc.Login(context.Background(), &entity.User{
		PhoneNumber: "081338442777",
		AuthData: entity.AuthData{
			Password:   "password",
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Internal, domain.ErrCreateTokenFailed))

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}
