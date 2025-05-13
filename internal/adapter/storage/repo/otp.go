package repo

import (
	"context"

	"go.bankyaya.org/app/backend/internal/adapter/storage/model"
	"go.bankyaya.org/app/backend/internal/domain/otp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OTPRepo struct {
	db *gorm.DB
}

func NewOTPRepo(db *gorm.DB) *OTPRepo {
	return &OTPRepo{
		db: db,
	}
}

func (o *OTPRepo) Save(ctx context.Context, otp *otp.OTP) error {
	m := &model.OTP{
		Code:       otp.Code,
		UserID:     otp.User.ID,
		Purpose:    otp.Purpose.String(),
		Channel:    otp.Channel.String(),
		VerifiedAt: otp.VerifiedAt,
		ExpiredAt:  otp.ExpiredAt,
	}
	res := o.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(m)
	return res.Error
}

func (o *OTPRepo) Get(ctx context.Context, id int) (*otp.OTP, error) {
	m := new(model.OTP)
	res := o.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		Find(m)
	if err := res.Error; err != nil {
		return nil, err
	}
	return &otp.OTP{
		ID:      m.ID,
		Code:    m.Code,
		Purpose: otp.NewPurpose(m.Purpose),
		Channel: otp.NewChannel(m.Channel),
		User: &otp.User{
			ID:    m.User.ID,
			Name:  m.User.FullName,
			Email: m.User.Email,
			Phone: m.User.PhoneNumber,
		},
		CreatedAt:  m.CreatedAt,
		ExpiredAt:  m.ExpiredAt,
		VerifiedAt: m.VerifiedAt,
	}, nil
}

func (o *OTPRepo) Update(ctx context.Context, otp *otp.OTP) error {
	res := o.db.WithContext(ctx).Updates(otp)
	return res.Error
}
