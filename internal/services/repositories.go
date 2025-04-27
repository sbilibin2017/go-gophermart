package services

import "context"

type ExistsRepository interface {
	Exists(ctx context.Context, filter map[string]any) (bool, error)
}

type FilterRepository interface {
	Filter(ctx context.Context, filter map[string]any, fields []string) (map[string]any, error)
}

type SaveRepository interface {
	Save(ctx context.Context, data map[string]any) error
}

type FilterILikeRepository interface {
	FilterILike(ctx context.Context, s string, fields []string) (map[string]any, error)
}

type ListRepository interface {
	List(ctx context.Context, filter map[string]any, fields []string, orderKey string, isOrderDesc bool) ([]map[string]any, error)
}

type PasswordHasher interface {
	Hash(password string) *string
}

type PasswordComparer interface {
	Compare(enteredPassword, storedPasswordHash string) error
}

type JWTGenerator interface {
	Generate(payload map[string]any) *string
}

type StructValidator interface {
	Struct(v any) error
}
