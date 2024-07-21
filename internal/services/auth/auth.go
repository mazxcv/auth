package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, UserID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAppNotFound        = errors.New("application not found")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
		log:          log,
	}
}

// Login checks if user with given credentials exists in system
//
// If user exists, but password is incorrect, returns error
// If user doesn't exist, return error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appId int64,
) (string, error) {
	const op = "auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to login user")
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("User not found: " + err.Error())
			return "", fmt.Errorf("%s, %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get User: " + err.Error())

		return "", fmt.Errorf("#{op}: #{err}")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("Invalid credentials: " + err.Error())
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, int(appId))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("User loggen successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate: " + err.Error())
		return "", fmt.Errorf("#{op}: #{err}")
	}

	return token, nil

}

// RegisterNewUser registers new user in the system and returns user ID
//
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash: " + err.Error())

		return 0, fmt.Errorf("#{op}: #{err}")
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("User already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user: " + err.Error())

		return 0, fmt.Errorf("#{op}: #{err}")
	}

	log.Info("new user registered")

	return id, nil
}

// IsAdmin check if user is admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, ErrAppNotFound)
		}
		return false, fmt.Errorf("#{op}: #{err}")
	}

	return isAdmin, nil
}
