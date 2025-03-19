package transfer

import (
	"context"
	"errors"

	transfer2 "go.bankyaya.org/app/backend/internal/core/transfer"
	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
	"go.bankyaya.org/app/backend/pkg/types"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type Usecase struct {
	va  *validation.Validator
	log *logger.Logger
	svc *transfer2.Service
}

func NewUsecase(va *validation.Validator, log *logger.Logger) *Usecase {
	return &Usecase{
		va:  va,
		log: log,
	}
}

func (uc *Usecase) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	sequence, err := uc.svc.Inquiry(ctx, &transfer2.Sequence{})
	if err != nil {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("Inquiry failed: %v", err)
		if errors.Is(err, transfer2.ErrEODInProgress) {
			return nil, status.Error(codes.Internal, "EOD process is running")
		}
		if errors.Is(err, transfer2.ErrSourceAccountInactive) {
			return nil, status.Error(codes.BadRequest, "Source account is inactive")
		}
		if errors.Is(err, transfer2.ErrDestinationAccountInactive) {
			return nil, status.Error(codes.BadRequest, "Source account is inactive")
		}
		return nil, status.Error(codes.Internal, "Inquiry failed")
	}

	return &InquiryResponse{
		SequenceNumber:     sequence.SeqNo,
		SourceAccount:      sequence.AccNoSrc,
		DestinationAccount: sequence.AccNoDest,
	}, nil
}

func (uc *Usecase) DoPayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	transaction, err := uc.svc.DoPayment(ctx, req.Sequence)
	if err != nil && errors.Is(err, transfer2.ErrSendEmailFailed) {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("Send email failed: %v", err)
	}
	if err != nil {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("Payment failed: %v", err)
		if errors.Is(err, transfer2.ErrEODInProgress) {
			return nil, status.Error(codes.Internal, "EOD process is running")
		}
		if errors.Is(err, transfer2.ErrInvalidSequenceNumber) {
			return nil, status.Error(codes.BadRequest, "Invalid sequence number")
		}
		if errors.Is(err, ctxt.ErrUserFromContext) {
			return nil, status.Error(codes.Unauthenticated, "Unauthorized user")
		}
		return nil, status.Error(codes.Internal, "Payment failed")
	}

	amount, err := types.ParseMoney(transaction.Amount)
	if err != nil {
		uc.log.Usecase("Transfer", "Inquiry").Errorf("ParseMoney error: %v", err)
		amount = 0
	}

	return &PaymentResponse{
		Amount: amount,
	}, nil
}
