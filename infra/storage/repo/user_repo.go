package repo

import (
	"context"

	"go.bankyaya.org/app/backend/domain/user"
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
	u := new(user.User)
	res := r.db.WithContext(ctx).
		Preload("AuthData").
		Where(`"PHONE_NUMBER" = ?`, phoneNumber).
		First(u)
	if err := res.Error; err != nil {
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
		return nil, err
	}
	return device, nil
}
