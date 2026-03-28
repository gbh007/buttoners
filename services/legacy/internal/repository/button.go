package repository

import (
	"context"

	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"github.com/gbh007/buttoners/services/legacy/internal/repository/gen"
	"gorm.io/gorm"
)

func (repo Repository) SetButton(ctx context.Context, b domain.Button) error {
	res := repo.db.Save(&b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) GetButton(ctx context.Context, b *domain.Button) error {
	res := repo.db.Where(b).First(b)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo Repository) ButtonsByUser(ctx context.Context, userID int) ([]domain.Button, error) {
	buttons, err := gorm.G[domain.Button](repo.db).
		Where(gen.Button.UserID.Eq(userID)).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	return buttons, nil
}

func (repo Repository) ButtonsTotalByUser(ctx context.Context, userID int) (int, error) {
	var c int

	err := gorm.G[domain.Button](repo.db).
		Where(gen.Button.UserID.Eq(userID)).
		Select("SUM("+gen.Button.Count.Column().Name+")").
		Scan(ctx, &c)
	if err != nil {
		return 0, err
	}

	return c, nil
}
