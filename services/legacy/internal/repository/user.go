package repository

import (
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"context"
)

func (repo Repository) SetUser(ctx context.Context, u domain.User) error {
	res := repo.db.Save(&u)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) GetUser(ctx context.Context, u *domain.User) error {
	res := repo.db.Model(u).First(u)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) GetUserByToken(ctx context.Context, u *domain.User) error {
	res := repo.db.Model(&domain.User{}).Where("token = ?", u.Token).First(u)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
