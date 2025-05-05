package otp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

func TestSendSuccess(t *testing.T) {
	var (
		repoMock      = NewMockRepository(t)
		generatorMock = NewMockGenerator(t)
		senderMock    = NewMockSender(t)
		svc           = NewService(logger.New(), repoMock, generatorMock, senderMock)
		ctx           = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
			Phone: "081234567890",
		})
	)

	generatorMock.EXPECT().Generate(otpLength).
		Return("123456", nil)

	senderMock.EXPECT().Send(mock.Anything, mock.Anything).
		Return(nil)

	repoMock.EXPECT().Save(mock.Anything, mock.Anything).
		Return(nil)

	res, err := svc.Send(ctx, PurposeLogin, ChannelEmail)

	assert.NoError(t, err)
	assert.Equal(t, "123456", res.Code)
	assert.Equal(t, PurposeLogin, res.Purpose)
	assert.Equal(t, ChannelEmail, res.Channel)

	generatorMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	senderMock.AssertExpectations(t)
}

func TestSend_GetUserFromContextFailed(t *testing.T) {
	var (
		repoMock      = NewMockRepository(t)
		generatorMock = NewMockGenerator(t)
		senderMock    = NewMockSender(t)
		svc           = NewService(logger.New(), repoMock, generatorMock, senderMock)
		ctx           = context.Background()
	)

	res, err := svc.Send(ctx, PurposeLogin, ChannelEmail)

	assert.Nil(t, res)
	assert.Equal(t, status.Error(codes.Unauthenticated, ErrUnauthenticatedUser), err)

	generatorMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	senderMock.AssertExpectations(t)
}

func TestSend_GenerateOTPFailed(t *testing.T) {
	var (
		repoMock      = NewMockRepository(t)
		generatorMock = NewMockGenerator(t)
		senderMock    = NewMockSender(t)
		svc           = NewService(logger.New(), repoMock, generatorMock, senderMock)
		ctx           = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
			Phone: "081234567890",
		})
	)

	generatorMock.EXPECT().Generate(otpLength).
		Return("", errors.New("failed to generate otp"))

	res, err := svc.Send(ctx, PurposeLogin, ChannelEmail)

	assert.Nil(t, res)
	assert.Equal(t, status.Error(codes.Internal, ErrGeneral), err)

	generatorMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	senderMock.AssertExpectations(t)
}

func TestVerifySuccess(t *testing.T) {
	var (
		repoMock      = NewMockRepository(t)
		generatorMock = NewMockGenerator(t)
		senderMock    = NewMockSender(t)
		svc           = NewService(logger.New(), repoMock, generatorMock, senderMock)
		ctx           = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
			Phone: "081234567890",
		})
	)

	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Minute * 5)

	repoMock.EXPECT().Get(mock.Anything, mock.Anything).
		Return(&OTP{
			ID:      1,
			Code:    "123456",
			Purpose: PurposeLogin,
			Channel: ChannelEmail,
			User: &User{
				ID:    123,
				Name:  "Olivia Rodrigo",
				Email: "olivia@gmail.com",
				Phone: "081234567890",
			},
			CreatedAt: createdAt,
			ExpiredAt: expiredAt,
		}, nil)
	repoMock.EXPECT().Update(mock.Anything, mock.Anything).
		Return(nil)

	err := svc.Verify(ctx, &OTP{
		ID:        1,
		Code:      "123456",
		Purpose:   PurposeLogin,
		Channel:   ChannelEmail,
		CreatedAt: createdAt,
		ExpiredAt: expiredAt,
	})

	assert.NoError(t, err)

	generatorMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	senderMock.AssertExpectations(t)
}
