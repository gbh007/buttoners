package storage

import (
	"context"
	"fmt"
	"time"
)

type TaskResult struct {
	UserID     int64     `db:"user_id"`
	Chance     int64     `db:"chance"`
	Duration   int64     `db:"duration"`
	Result     int64     `db:"result"`
	ResultText string    `db:"result_text"`
	ErrorText  string    `db:"error_text"`
	StartTime  time.Time `db:"start_time"`
	EndTime    time.Time `db:"end_time"`
}

func (db *Database) InsertTaskResult(ctx context.Context, t *TaskResult) error {
	t.StartTime = t.StartTime.UTC()
	t.EndTime = t.EndTime.UTC()

	_, err := db.db.NamedExecContext(ctx, `INSERT INTO task_results (
        user_id,
        chance,
        duration,
        result,
        result_text,
        error_text,
        start_time,
        end_time
) VALUES (
        :user_id,
        :chance,
        :duration,
        :result,
        :result_text,
        :error_text,
        :start_time,
        :end_time
);`, t)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}
