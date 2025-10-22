package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"auth/storage"

	"golang.org/x/crypto/bcrypt"
)

type AuthApp struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *AuthApp {
	return &AuthApp{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		tokenTTL:    tokenTTL, // Время жизни возвращаемых токенов
	}
}

func (a *AuthApp) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "Auth.RegisterNewUser"

	// Создаём локальный объект логгера с доп. полями, содержащими полезную инфу
	// о текущем вызове функции
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	// Генерируем хэш и соль для пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем пользователя в БД
	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (a *AuthApp) Login(
	ctx context.Context,
	email string,
	password string, // пароль в чистом виде
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	// Достаём пользователя из БД
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	// Создаём токен авторизации
	token, err := NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
