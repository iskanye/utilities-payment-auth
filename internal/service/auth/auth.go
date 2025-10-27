package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/iskanye/utilities-payment-auth/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-auth/internal/storage"
	"github.com/iskanye/utilities-payment/pkg/logger"
	"github.com/iskanye/utilities-payment/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	secret      string
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
		isAdmin bool,
	) (uid int64, err error)
}

type UserProvider interface {
	User(
		ctx context.Context,
		email string,
	) (models.User, error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (bool, error)
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	secret string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		secret:      secret,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
) (string, int64, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", logger.Err(err))

			return "", 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", logger.Err(err))

		return "", 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", logger.Err(err))

		return "", 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, a.secret, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", logger.Err(err))

		return "", 0, fmt.Errorf("%s: %w", op, err)
	}

	return token, user.ID, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) Register(
	ctx context.Context,
	email string,
	pass string,
) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", logger.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash, false)
	if err != nil {
		log.Error("failed to save user", logger.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) Validate(
	ctx context.Context,
	token string,
) (bool, error) {
	const op = "Auth.Validate"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("validating user")

	isValid, err := jwt.Validate(token, a.secret)
	if err != nil {
		log.Error("failed to validate user", logger.Err(err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isValid, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
