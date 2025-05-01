package repo

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/adapter/storage/model"
	"go.bankyaya.org/app/backend/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*user.User, error) {
	u := new(model.User)
	res := r.db.WithContext(ctx).
		Preload("AuthData").
		Preload("AuthData.Device").
		Where(`"PHONE_NUMBER" = ?`, phoneNumber).
		First(u)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &user.User{
		ID:             u.ID,
		AccountNumber:  u.AccountNumber,
		FullName:       u.FullName,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		IdentityNumber: u.IdentityNo,
		Device: &user.Device{
			FirebaseId:    u.AuthData.FirebaseId,
			DeviceId:      u.AuthData.DeviceId,
			IsBlacklisted: u.AuthData.Device.IsBlacklisted(),
		},
	}, nil
}
