package button

import (
	"context"
	"errors"
	"io"
	"slices"
	"time"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/services/legacy/internal/domain"
	"github.com/gbh007/buttoners/services/legacy/internal/metrics"
	"github.com/gbh007/buttoners/services/legacy/internal/repository"
	"github.com/segmentio/ksuid"

	"gorm.io/gorm"
)

type Service struct {
	repo            *repository.Repository
	kafkaTaskClient *kafka.Producer[dto.KafkaTaskData]
}

func New(repo *repository.Repository, kafkaTaskClient *kafka.Producer[dto.KafkaTaskData]) *Service {
	return &Service{
		repo:            repo,
		kafkaTaskClient: kafkaTaskClient,
	}
}

func (s *Service) PressButton(ctx context.Context, user domain.User) error {
	key := ksuid.New().String()

	err := s.kafkaTaskClient.Write(ctx, key, dto.KafkaTaskData{
		UserID:   int64(user.ID),
		Chance:   50,
		Duration: 1,
	})
	if err != nil {
		return err
	}

	metrics.RecordButtonPress(user.Name)

	return nil
}

func (s *Service) ConsumePressButton(ctx context.Context, userID int) error {
	y, m, d := time.Now().Date()

	b := domain.Button{
		UserID: userID,
		Year:   y,
		Month:  int(m),
		Day:    d,
	}

	err := s.repo.GetButton(ctx, &b)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	b.Count++

	err = s.repo.SetButton(ctx, b)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Buttons(ctx context.Context, user domain.User) ([]domain.Button, error) {
	buttons, err := s.repo.ButtonsByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(buttons, func(a, b domain.Button) int {
		if a.Year != b.Year {
			return b.Year - a.Year
		}

		if a.Month != b.Month {
			return b.Month - a.Month
		}

		if a.Day != b.Day {
			return b.Day - a.Day
		}

		return 0
	})

	if len(buttons) > 7 {
		buttons = buttons[:7]
	}

	for i := range buttons {
		buttons[i].UpdateText()
	}

	return buttons, nil
}

func (s *Service) ButtonBadge(ctx context.Context, w io.Writer, userID int) error {
	c, err := s.repo.ButtonsTotalByUser(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return renderBadgeTemplate(w, c)
}
