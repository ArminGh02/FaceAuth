package repository

import (
	"context"
	"errors"

	"github.com/ArminGh02/go-auth-system/internal/model"
)

type User interface {
	Insert(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByNationalID(ctx context.Context, nationalID string) (*model.User, error)
	Close(ctx context.Context) error
}

var ErrNotFound = errors.New("user not found")
