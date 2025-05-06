package otp

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/pkgerror"
)

const (
	domainName = "otp"
	otpLength  = 6
	otpExpiry  = 5 * time.Minute
)

type Service struct {
	log       *logger.Logger
	repo      Repository
	generator Generator
	sender    Sender
}

func NewService(
	log *logger.Logger,
	repo Repository,
	generator Generator,
	sender Sender,
) *Service {
	return &Service{
		log:       log,
		repo:      repo,
		generator: generator,
		sender:    sender,
	}
}

func (s *Service) Send(ctx context.Context, purpose Purpose, channel Channel) (*OTP, error) {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.DomainUsecase(domainName, "Send").Error(ctxt.ErrUserFromContext)
		return nil, pkgerror.New(codes.Unauthenticated, ErrUnauthenticatedUser)
	}

	code, err := s.generator.Generate(otpLength)
	if err != nil {
		s.log.DomainUsecase(domainName, "Send").Error(err)
		return nil, pkgerror.New(codes.Internal, ErrGeneral)
	}

	createdAt := time.Now()
	expiredAt := createdAt.Add(otpExpiry)

	otp := &OTP{
		Code:    code,
		Purpose: purpose,
		Channel: channel,
		User: &User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
		CreatedAt: createdAt,
		ExpiredAt: expiredAt,
	}

	err = s.sender.Send(ctx, otp)
	if err != nil {
		s.log.DomainUsecase(domainName, "Send").Error(err)
		return nil, pkgerror.New(codes.Internal, ErrGeneral)
	}

	err = s.repo.Save(ctx, otp)
	if err != nil {
		s.log.DomainUsecase(domainName, "Send").Error(err)
		return nil, pkgerror.New(codes.Internal, ErrGeneral)
	}

	return otp, nil
}

func (s *Service) Verify(ctx context.Context, in *OTP) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.DomainUsecase(domainName, "Verify").Error(ctxt.ErrUserFromContext)
		return pkgerror.New(codes.Unauthenticated, ErrUnauthenticatedUser)
	}

	in.User = NewUser(user.ID, user.Name, user.Email, user.Phone)

	otp, err := s.repo.Get(ctx, in.ID)
	if err != nil {
		s.log.DomainUsecase(domainName, "Verify").Error(err)
		return pkgerror.New(codes.Internal, ErrGeneral)
	}
	if !otp.Equal(in) {
		s.log.DomainUsecase(domainName, "Verify").Error(ErrInvalidOTP)
		return pkgerror.New(codes.BadRequest, ErrInvalidOTP)
	}
	if otp.IsVerified() {
		s.log.DomainUsecase(domainName, "Verify").Error(ErrOTPAlreadyUsed)
		return pkgerror.New(codes.BadRequest, ErrOTPAlreadyUsed)
	}
	if otp.IsExpired(time.Now()) {
		s.log.DomainUsecase(domainName, "Verify").Error(ErrOTPExpired)
		return pkgerror.New(codes.BadRequest, ErrOTPExpired)
	}

	otp.VerifiedAt = time.Now()

	err = s.repo.Update(ctx, otp)
	if err != nil {
		s.log.DomainUsecase(domainName, "Verify").Error(err)
		return pkgerror.New(codes.Internal, ErrGeneral)
	}

	return nil
}
