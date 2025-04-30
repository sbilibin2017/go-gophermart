package registries

import "github.com/sbilibin2017/go-gophermart/internal/types"

type HTTPErrorRegistry struct {
	errors map[string]*types.HTTPError
}

func NewHTTPErrorRegistry() *HTTPErrorRegistry {
	return &HTTPErrorRegistry{
		errors: make(map[string]*types.HTTPError),
	}
}

func (r *HTTPErrorRegistry) Register(err *types.HTTPError) {
	r.errors[err.Error.Error()] = err
}

func (r *HTTPErrorRegistry) Get(err error) *types.HTTPError {
	if registeredError, exists := r.errors[err.Error()]; exists {
		return registeredError
	}
	return nil
}
