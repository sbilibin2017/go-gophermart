package helpers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
