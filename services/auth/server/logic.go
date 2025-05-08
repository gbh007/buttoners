package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gbh007/buttoners/services/auth/internal/storage"
)

var (
	// Неверный логин/пароль
	ErrLoginOrPasswordIncorrect = errors.New("login or password incorrect")
	// Сессия не найдена
	ErrSessionNotFound = errors.New("session not found")
)

// createUser - создает нового пользователя
func (s *authServer) createUser(ctx context.Context, login, password string) (int64, error) {
	salt := randomSHA256String()
	login = strings.ToLower(login)

	id, err := s.db.CreateUser(ctx, &storage.User{
		Login:    login,
		Password: saltPassword(password, salt),
		Salt:     salt,
	})
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

// createSession - создает новую сессию пользователя
func (s *authServer) createSession(ctx context.Context, login, password string) (string, error) {
	user, err := s.checkUser(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	token := randomSHA256String()

	err = s.db.CreateSession(ctx, &storage.Session{
		Token:  token,
		UserID: user.ID,
	})
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return token, nil
}

// deleteSession - удаляет сессию пользователя
func (s *authServer) deleteSession(ctx context.Context, token string) error {
	err := s.db.DeleteSessionByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

// checkUser - проверяет данные пользователя
func (s *authServer) checkUser(ctx context.Context, login, password string) (*storage.User, error) {
	login = strings.ToLower(login)

	user, err := s.db.GetUserByLogin(ctx, login)

	// Такого пользователя не существует
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLoginOrPasswordIncorrect
	}

	if err != nil {
		return nil, err
	}

	// Проверка пароля
	if saltPassword(password, user.Salt) != user.Password {
		return nil, ErrLoginOrPasswordIncorrect
	}

	return user, nil
}

// getUser - возвращает данные пользователя по токену
func (s *authServer) getUser(ctx context.Context, token string) (*storage.User, error) {
	session, err := s.db.GetSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	user, err := s.db.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
