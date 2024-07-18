package auth

import (
	"context"
	"log/slog"
	"sso/internal/domain/models"
	"time"
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
	User(ctx context.Context, UserID int64) (models.User, error)
	IsAdmin(ctx context.Context, UserID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

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
	addID int,
) (string, error) {
	panic("not implemented")
}

// RegisterNewUser registers new user in the system and returns user ID
//
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	panic("not implemented")
}

// IsAdmin check if user is admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	addID int,
) (bool, error) {
	panic("not implemented")
}
