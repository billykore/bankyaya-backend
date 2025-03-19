package repo

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/core/user"
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

func (r *UserRepo) GetUserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*user.User, error) {
	u := new(user.User)
	res := r.db.WithContext(ctx).
		Preload("AuthData").
		Preload("AuthData.Device").
		Where(`"PHONE_NUMBER" = ?`, phoneNumber).
		First(u)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetDeviceById(ctx context.Context, deviceId string) (*user.Device, error) {
	device := new(user.Device)
	res := r.db.WithContext(ctx).
		Preload("Blacklist").
		Where(`"DEVICE_ID" = ?`, deviceId).
		First(device)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}
