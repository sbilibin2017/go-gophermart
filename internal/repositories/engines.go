package repositories

import "context"

type Querier interface {
	Query(ctx context.Context, query string, dest any, args any) error
}

type Executor interface {
	Exec(ctx context.Context, query string, args any) error
}
