package repo

import (
	"context"

	"go.bankyaya.org/app/backend/internal/adapter/storage/model"
	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"gorm.io/gorm"
)

const intrabankTransactionType = "internal_transfer"

type IntrabankRepo struct {
	db *gorm.DB
}

func NewIntrabankRepo(db *gorm.DB) *IntrabankRepo {
	return &IntrabankRepo{
		db: db,
	}
}

func (repo *IntrabankRepo) GetTransactionLimit(ctx context.Context) (*intrabank.Limits, error) {
	txLimit := new(model.TransactionMethod)
	res := repo.db.WithContext(ctx).
		Select(`"TRANSACTION_MIN_LIMIT", "TRANSACTION_LIMIT", "DAILY_LIMIT"`).
		Where(`"TYPE" = ?`, intrabankTransactionType).
		Find(txLimit)
	if err := res.Error; err != nil {
		return nil, err
	}
	minAmount, err := intrabank.ParseMoney(txLimit.TransactionMinLimit)
	if err != nil {
		return nil, err
	}
	return &intrabank.Limits{
		MinAmount:      minAmount,
		MaxAmount:      intrabank.Money(txLimit.TransactionLimit),
		MaxDailyAmount: intrabank.Money(txLimit.DailyLimit),
	}, nil
}

func (repo *IntrabankRepo) InsertSequence(ctx context.Context, seq *intrabank.Sequence) error {
	res := repo.db.WithContext(ctx).Create(seq)
	return res.Error
}

func (repo *IntrabankRepo) GetSequence(ctx context.Context, sequenceNumber string) (*intrabank.Sequence, error) {
	seq := new(intrabank.Sequence)
	res := repo.db.WithContext(ctx).
		Where(`"SEQ_NO" = ?`, sequenceNumber).
		First(seq)
	if err := res.Error; err != nil {
		return nil, err
	}
	return seq, nil
}

func (repo *IntrabankRepo) InsertTransaction(ctx context.Context, transaction *intrabank.Transaction) error {
	res := repo.db.WithContext(ctx).Create(transaction)
	return res.Error
}
