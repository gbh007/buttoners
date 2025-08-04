package repository

import (
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"context"
)

func (repo Repository) SetButton(ctx context.Context, b domain.Button) error {
	res := repo.db.Save(&b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) GetButton(ctx context.Context, b *domain.Button) error {
	res := repo.db.Model(b).Take(b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) ButtonsByUser(ctx context.Context, userID int) ([]domain.Button, error) {
	buttons := make([]domain.Button, 0)

	res := repo.db.Model(&domain.Button{}).Where("user_id = ?", userID).Find(&buttons)
	if res.Error != nil {
		return nil, res.Error
	}

	return buttons, nil
}

func (repo Repository) ButtonsTotalByUser(ctx context.Context, userID int) (int, error) {
	var c int

	res := repo.db.Model(&domain.Button{}).Where("user_id = ?", userID).Select("sum(count)").Scan(&c)
	if res.Error != nil {
		return 0, res.Error
	}

	return c, nil
}
