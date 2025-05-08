package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// User - данные пользователя
type User struct {
	// Ид в базе
	ID int64 `db:"id"`
	// Логин пользователя
	Login string `db:"login"`
	// Засоленный пароль пользователя
	Password string `db:"password"`
	// Соль для пользователя
	Salt string `db:"salt"`
	// Время создания учетной записи
	Created time.Time `db:"created"`
	// Время последнего обновления учетной записи
	Updated sql.NullTime `db:"updated"`
}

// CreateUser - создает нового пользователя в базе.
// Поля ID, Updated игнорируются, поле Created заменяется
func (d *Database) CreateUser(ctx context.Context, user *User) (int64, error) {
	var id int64

	user.Created = time.Now().UTC()

	query, args, err := d.db.BindNamed(
		`INSERT INTO users (login, password, salt, created) VALUES (:login, :password, :salt, :created) RETURNING id;`,
		user,
	)
	if err != nil {
		return 0, fmt.Errorf("%w; %w", errDatabase, err)
	}

	err = d.db.GetContext(ctx, &id, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%w; %w", errDatabase, err)
	}

	return id, nil
}

// GetUserByID - возвращает пользователя по ИД
func (d *Database) GetUserByID(ctx context.Context, id int64) (*User, error) {
	user := new(User)

	err := d.db.GetContext(ctx, user, `SELECT * FROM users WHERE id = ? LIMIT 1;`, id)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", errDatabase, err)
	}

	return user, nil
}

// GetUserByID - возвращает пользователя по логину
func (d *Database) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	user := new(User)

	err := d.db.GetContext(ctx, user, `SELECT * FROM users WHERE login = ? LIMIT 1;`, login)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", errDatabase, err)
	}

	return user, nil
}

// GetUsers - возвращает всех пользователей
func (d *Database) GetUsers(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0)

	err := d.db.SelectContext(ctx, &users, `SELECT * FROM users ORDER BY ID;`)
	if err != nil {
		return nil, fmt.Errorf("%w; %w", errDatabase, err)
	}

	return users, nil
}
