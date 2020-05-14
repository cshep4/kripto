package service

import (
	"context"
	"errors"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
)

type (
	Storer interface {
		Store(ctx context.Context, user model.User) error
		GetByEmail(ctx context.Context, email string) (*model.User, error)
	}

	service struct {
		store Storer
	}
)

func New(store Storer) (*service, error) {
	if store == nil {
		return nil, errors.New("store_is_nil")
	}

	return &service{
		store: store,
	}, nil
}

func (s *service) Register(ctx context.Context, user model.User) error {
	panic("implement me")
}
