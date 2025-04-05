package repo

import (
	"context"
	"errors"

	pkgerrors "go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
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

func (r *UserRepo) GetUserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.User, error) {
	u := new(entity.User)
	res := r.db.WithContext(ctx).
		Preload("AuthData").
		Preload("AuthData.Device").
		Where(`"PHONE_NUMBER" = ?`, phoneNumber).
		First(u)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetDeviceById(ctx context.Context, deviceId string) (*entity.Device, error) {
	device := new(entity.Device)
	res := r.db.WithContext(ctx).
		Preload("Blacklist").
		Where(`"DEVICE_ID" = ?`, deviceId).
		First(device)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}
