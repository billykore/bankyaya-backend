package user

import (
	"context"
	"strings"

	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/data"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/security/password"
	"go.bankyaya.org/app/backend/pkg/security/token"
	"go.bankyaya.org/app/backend/pkg/status"
)

const (
	messageLoginFailed         = "Login failed"
	messageDeviceIsBlacklisted = "Device is blacklisted"
	messageUserNotFound        = "User not found"
	messageInvalidPassword     = "Invalid password"
	messageInvalidDevice       = "Invalid device"
)

// Service handles user related process.
type Service struct {
	log  *logger.Logger
	repo Repository
}

func NewService(log *logger.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	device, err := s.repo.GetDeviceById(ctx, req.DeviceId)
	if err != nil {
		s.log.DomainUsecase("user", "Login").Errorf("GetDeviceById: %v", err)
		return nil, status.Error(codes.Internal, messageLoginFailed)
	}
	if device.IsBlacklisted() {
		s.log.DomainUsecase("user", "Login").Errorf("deviceId (%s) is blacklisted", req.DeviceId)
		return nil, status.Error(codes.Forbidden, messageDeviceIsBlacklisted)
	}

	user, err := s.repo.GetUserByPhoneNumber(ctx, req.Phone)
	if err != nil {
		s.log.DomainUsecase("user", "Login").Errorf("GetUserByPhoneNumber: %v", err)
		return nil, status.Errorf(codes.NotFound, "%v: %v", messageLoginFailed, messageUserNotFound)
	}

	err = password.Verify(user.AuthData.Password, req.Password)
	if err != nil {
		s.log.DomainUsecase("user", "Login").Errorf("Verify: %v", err)
		return nil, status.Errorf(codes.Forbidden, "%v: %v", messageLoginFailed, messageInvalidPassword)
	}

	if strings.Compare(req.DeviceId, user.AuthData.DeviceId) != 0 && strings.Compare(req.FirebaseId, user.AuthData.FirebaseId) != 0 {
		s.log.DomainUsecase("user", "Login").Error(ErrInvalidDevice)
		return nil, status.Errorf(codes.Forbidden, "%v: %v", messageLoginFailed, messageInvalidDevice)
	}

	accessToken, err := token.New(data.User{
		Id:       user.ID,
		CIF:      user.CIF,
		FullName: user.FullName,
	})
	if err != nil {
		s.log.DomainUsecase("user", "Login").Errorf("New token: %v", err)
		return nil, status.Errorf(codes.Internal, messageLoginFailed)
	}

	return &LoginResponse{
		Token:       accessToken.AccessToken,
		ExpiredTime: accessToken.ExpiredTime,
	}, nil
}
