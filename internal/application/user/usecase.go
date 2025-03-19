package user

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/core/user"
	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type Usecase struct {
	va  *validation.Validator
	log *logger.Logger
	svc *user.Service
}

func NewUsecase(va *validation.Validator, log *logger.Logger, svc *user.Service) *Usecase {
	return &Usecase{
		va:  va,
		log: log,
		svc: svc,
	}
}

func (uc *Usecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("user", "Login").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	token, err := uc.svc.Login(ctx, &user.User{
		PhoneNumber: req.Phone,
		AuthData: user.AuthData{
			Password:   req.Password,
			DeviceId:   req.DeviceId,
			FirebaseId: req.FirebaseId,
		},
	})
	if err != nil {
		uc.log.Usecase("user", "Login").Errorf("Login failed: %v", err)
		if errors.Is(err, user.ErrDeviceIsBlacklisted) {
			return nil, status.Error(codes.BadRequest, "Device is blacklisted")
		}
		if errors.Is(err, user.ErrInvalidDevice) {
			return nil, status.Error(codes.BadRequest, "Device is invalid")
		}
		if errors.Is(err, user.ErrInvalidPassword) {
			return nil, status.Error(codes.BadRequest, "Invalid username or password")
		}
		return nil, status.Error(codes.Internal, "Login failed")
	}

	uc.log.Usecase("user", "Login").Infof("User (%v) login successfully", req.Phone)
	return &LoginResponse{
		Token:       token.AccessToken,
		ExpiredTime: token.ExpiredTime,
	}, nil
}
