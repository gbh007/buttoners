package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserLog struct {
	RequestID    string         `db:"request_id"`
	Addr         string         `db:"addr"`
	UserID       sql.NullInt64  `db:"user_id"`
	SessionToken sql.NullString `db:"session_token"`
	Action       string         `db:"action"`
	Chance       sql.NullInt64  `db:"chance"`
	Duration     sql.NullInt64  `db:"duration"`
	RequestTime  time.Time      `db:"request_time"`
}

func (db *Database) InsertUserLog(ctx context.Context, ul *UserLog) error {
	ul.RequestTime = ul.RequestTime.UTC()

	_, err := db.db.NamedExecContext(ctx, `INSERT INTO user_logs (
        request_id,
        addr,
        user_id,
        session_token,
        action,
        chance,
        duration,
        request_time
) VALUES (
        :request_id,
        :addr,
        :user_id,
        :session_token,
        :action,
        :chance,
        :duration,
        :request_time
);`, ul)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}

func (db *Database) SelectUserLogByUserID(ctx context.Context, userID int64) ([]*UserLog, error) {
	logs := make([]*UserLog, 0)

	err := db.db.SelectContext(ctx, &logs, `SELECT * FROM user_logs WHERE user_id = ?;`, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	return logs, nil
}

func (db *Database) SelectCompressedUserLogByUserID(ctx context.Context, userID int64) (int64, time.Time, error) {
	var (
		count int64
		last  time.Time
	)

	row := db.db.QueryRowContext(
		ctx,
		`SELECT COUNT(request_id), MAX(request_time) FROM user_logs WHERE user_id = ? GROUP BY user_id;`,
		userID,
	)

	err := row.Scan(&count, &last)
	if err != nil {
		return count, last, fmt.Errorf("%w: %w", errDatabase, err)
	}

	return count, last, nil
}
