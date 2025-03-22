package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/pkg/entity"
	repomock "go.bankyaya.org/app/backend/pkg/interface/repository/mocks"
	securitymock "go.bankyaya.org/app/backend/pkg/interface/security/mocks"
)

func TestLoginSuccess(t *testing.T) {
	mockRepo := repomock.NewUserRepository(t)
	mockRepo.On("GetUserDataByPhoneNumber", mock.Anything, "081338442777").
		Return(&entity.User{
			AuthData: entity.AuthData{
				FirebaseId: "123",
				DeviceId:   "456",
			},
		}, nil)

	mockHasher := securitymock.NewPasswordHasher(t)
	mockHasher.On("Compare", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(true)

	mockTokenSvc := securitymock.NewTokenService(t)
	mockTokenSvc.On("Create", mock.Anything, 15*time.Minute).
		Return(entity.Token{
			AccessToken: "example-token-123",
			ExpiredTime: time.Now().Add(15 * time.Minute).Unix(),
		}, nil)

	type args struct {
		ctx  context.Context
		user *entity.User
	}

	type want struct {
		token *entity.Token
		err   error
	}

	type test struct {
		caseName string
		args     args
		want     want
	}

	svc := NewUser(mockRepo, mockHasher, mockTokenSvc)

	for _, tt := range []test{
		{
			caseName: "success login",
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					PhoneNumber: "081338442777",
					AuthData: entity.AuthData{
						FirebaseId: "123",
						DeviceId:   "456",
					},
				},
			},
			want: want{
				token: &entity.Token{
					AccessToken: "example-token-123",
					ExpiredTime: time.Now().Add(15 * time.Minute).Unix(),
				},
				err: nil,
			},
		},
	} {
		t.Run(tt.caseName, func(t *testing.T) {
			token, err := svc.Login(tt.args.ctx, tt.args.user)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.token, token)
		})
	}

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockTokenSvc.AssertExpectations(t)
}
