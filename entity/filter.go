package entity

import "context"

type BaseFilterRepository[T any] interface {
	GetPaginated(ctx context.Context, filter *Filter) ([]T, int, error)
	GetFiltered(ctx context.Context, filter *Filter) ([]T, int, error)
}

type Filter struct {
	Page   int
	Limit  int
	Search string
	Sort   []string
	Filter []FilterOptions
}

type FilterOptions struct {
	Name     string
	Operator string
	Value    []string
}
