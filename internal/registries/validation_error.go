package registries

import (
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type ValidationErrorRegistry struct {
	errors map[string]*types.ValidationWithStatusCode
}

func NewValidationErrorRegistry() *ValidationErrorRegistry {
	return &ValidationErrorRegistry{
		errors: make(map[string]*types.ValidationWithStatusCode),
	}
}

func (h *ValidationErrorRegistry) Register(
	err types.ValidationError,
	statusCode int,
) {
	h.errors[err.Tag] = types.NewValidationWithStatusCode(err, statusCode)
}

func (h *ValidationErrorRegistry) Get(err error) *types.ValidationWithStatusCode {
	valErr := types.NewValidationError(err)
	if valErr == nil {
		return nil
	}
	if regErr, found := h.errors[valErr.Tag]; found {
		return regErr

	}
	return nil
}
