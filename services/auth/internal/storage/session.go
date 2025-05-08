package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Session - данные сессии пользователя
type Session struct {
	// Токен сессии
	Token string `db:"token"`
	// Ид пользователя в базе
	UserID int64 `db:"user_id"`
	// Время создания сессии
	Created time.Time `db:"created"`
	// Время последнего использования (обращения) сессии
	Used sql.NullTime `db:"used"`
	// Время последнего обновления сессии
	Updated sql.NullTime `db:"updated"`
}

// CreateSession - создает в базе запись о сессии пользователя
// Поля IsClosed, Used, Updated игнорируются, поле Created заменяется
func (d *Database) CreateSession(ctx context.Context, session *Session) error {
	session.Created = time.Now().UTC()

	_, err := d.db.NamedExecContext(
		ctx,
		`INSERT INTO sessions(token, user_id, created) VALUES (:token, :user_id, :created);`,
		session,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}

// GetSessionByToken - возвращает сессию пользователя по токену, не проверяет ее на закрытие.
func (d *Database) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	session := new(Session)

	err := d.db.GetContext(ctx, session, `SELECT * FROM sessions WHERE token = ? LIMIT 1;`, token)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabase, err)
	}

	return session, nil
}

// UpdateSessionUsedTime - обновляет время использования сессии.
func (d *Database) UpdateSessionUsedTime(ctx context.Context, token string) error {
	_, err := d.db.ExecContext(ctx, `UPDATE sessions SET used = ? WHERE token = ?;`, time.Now().UTC(), token)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}

// DeleteSessionByToken - удаляет сессию по токену.
func (d *Database) DeleteSessionByToken(ctx context.Context, token string) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = ?;`, token)
	if err != nil {
		return fmt.Errorf("%w: %w", errDatabase, err)
	}

	return nil
}
