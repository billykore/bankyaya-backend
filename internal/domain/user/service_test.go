package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

func TestSuccessLogin(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338442777").
		Return(&User{
			Password: "password",
			Device: &Device{
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.EXPECT().Compare("password", "password").
		Return(true)

	tokenSvcMock.EXPECT().Create(mock.Anything, 15*time.Minute).
		Return(&Token{
			AccessToken: "example-token-123",
			ExpiresAt:   time.Now().Add(15 * time.Minute),
		}, nil)

	token, err := svc.Login(context.Background(), &User{
		Password:    "password",
		PhoneNumber: "081338442777",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, "example-token-123", token.AccessToken)

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}

func TestLoginFailed_UserNotFound(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338000000").
		Return(nil, errors.New("user not found"))

	token, err := svc.Login(context.Background(), &User{
		PhoneNumber: "081338000000",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.NotFound, ErrUserNotFound))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_UserBlacklisted(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338000001").
		Return(&User{
			Device: &Device{
				FirebaseId:    "123",
				DeviceId:      "456",
				IsBlacklisted: true,
			},
		}, nil)

	token, err := svc.Login(context.Background(), &User{
		PhoneNumber: "081338000001",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Forbidden, ErrDeviceIsBlacklisted))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_InvalidDevice(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338000002").
		Return(&User{
			Device: &Device{
				FirebaseId: "321",
				DeviceId:   "654",
			},
		}, nil)

	token, err := svc.Login(context.Background(), &User{
		PhoneNumber: "081338000002",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Forbidden, ErrInvalidDevice))

	repoMock.AssertExpectations(t)
}

func TestLoginFailed_InvalidPassword(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338000003").
		Return(&User{
			Password: "password",
			Device: &Device{
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.EXPECT().Compare("invalid-password", "password").
		Return(false)

	token, err := svc.Login(context.Background(), &User{
		PhoneNumber: "081338000003",
		Password:    "invalid-password",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.BadRequest, ErrInvalidPassword))

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}

func TestLoginFailed_CreateTokenFailed(t *testing.T) {
	var (
		repoMock     = NewMockRepository(t)
		hasherMock   = NewMockPasswordHasher(t)
		tokenSvcMock = NewMockTokenService(t)
		svc          = NewService(logger.New(), repoMock, hasherMock, tokenSvcMock)
	)

	repoMock.EXPECT().GetUserByPhoneNumber(mock.Anything, "081338442777").
		Return(&User{
			Password: "password",
			Device: &Device{
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	hasherMock.EXPECT().Compare("password", "password").
		Return(true)

	tokenSvcMock.EXPECT().Create(mock.Anything, 15*time.Minute).
		Return(nil, errors.New("failed to create token"))

	token, err := svc.Login(context.Background(), &User{
		Password:    "password",
		PhoneNumber: "081338442777",
		Device: &Device{
			FirebaseId: "123",
			DeviceId:   "456",
		},
	})

	assert.Nil(t, token)
	assert.Equal(t, err, status.Error(codes.Internal, ErrCreateTokenFailed))

	repoMock.AssertExpectations(t)
	hasherMock.AssertExpectations(t)
	tokenSvcMock.AssertExpectations(t)
}
