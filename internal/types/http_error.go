package types

type HTTPError struct {
	Error      error
	StatusCode int
}
