package repo

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
	"gorm.io/gorm"
)

type TransferRepo struct {
	db *gorm.DB
}

func NewTransferRepo(db *gorm.DB) *TransferRepo {
	return &TransferRepo{
		db: db,
	}
}

func (repo *TransferRepo) InsertSequence(ctx context.Context, seq *entity.Sequence) error {
	res := repo.db.WithContext(ctx).Create(seq)
	return res.Error
}

func (repo *TransferRepo) GetSequence(ctx context.Context, sequenceNumber string) (*entity.Sequence, error) {
	seq := new(entity.Sequence)
	res := repo.db.WithContext(ctx).
		Where(`"SEQ_NO" = ?`, sequenceNumber).
		First(seq)
	if err := res.Error; err != nil {
		return nil, err
	}
	return seq, nil
}

func (repo *TransferRepo) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	user := new(entity.User)
	res := repo.db.WithContext(ctx).
		Where(`"ID" = ?`, id).
		First(user)
	if err := res.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *TransferRepo) InsertTransaction(ctx context.Context, transaction *entity.Transaction) error {
	res := repo.db.WithContext(ctx).Create(transaction)
	return res.Error
}
