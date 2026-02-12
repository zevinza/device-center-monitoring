package entity

import (
	"context"

	"github.com/google/uuid"
)

type BaseDomain[request any, response any] interface {
	GetAll(ctx context.Context) ([]response, error)
	GetByID(ctx context.Context, id *uuid.UUID) (*response, error)
	Create(ctx context.Context, request *request) (*response, error)
	Update(ctx context.Context, id *uuid.UUID, request *request) (*response, error)
	Delete(ctx context.Context, id *uuid.UUID) error
}

type BaseDomainPagination[response any] interface {
	GetPaginated(ctx context.Context, filter *Filter) ([]response, int64, error)
	GetFiltered(ctx context.Context, filter *Filter) ([]response, int64, error)
}
