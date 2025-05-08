package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Notification - данные уведомления пользователя
type Notification struct {
	// Ид в базе
	ID int64 `db:"id"`
	// Ид пользователя в базе
	UserID int64 `db:"user_id"`
	// Код уведомления (нажатие на кнопку, авторизация, и т.д.)
	Kind string `db:"kind"`
	// Уровень уведомления (информация, ошибка, и т.д.)
	Level string `db:"level"`
	// Заголовок уведомления
	Title string `db:"title"`
	// Тело уведомления
	Body sql.NullString `db:"body"`
	// Уведомление прочитано
	Read bool `db:"read"`
	// Время создания сессии
	Created time.Time `db:"created"`
}

func (d *Database) CreateNotification(ctx context.Context, session *Notification) error {
	session.Created = time.Now().UTC()

	_, err := d.db.NamedExecContext(ctx, `INSERT INTO notifications(
        user_id,
        kind,
        level,
        title,
        body,
        created
	) VALUES (
        :user_id,
        :kind,
        :level,
        :title,
        :body,
        :created
	);`, session)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}

func (d *Database) GetNotificationsByUserID(ctx context.Context, userID int64) ([]*Notification, error) {
	users := make([]*Notification, 0)

	err := d.db.SelectContext(ctx, &users, `SELECT * FROM notifications WHERE user_id = ? ORDER BY id;`, userID)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", errDatabase, err)
	}

	return users, nil
}

func (d *Database) MarkReadByID(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, "UPDATE notifications SET `read` = TRUE WHERE id = ?;", id)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}

func (d *Database) MarkReadByUserID(ctx context.Context, userID int64) error {
	_, err := d.db.ExecContext(ctx, "UPDATE notifications SET `read` = TRUE WHERE user_id = ?;", userID)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}
