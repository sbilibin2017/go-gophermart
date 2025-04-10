package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccrualRoutes(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		url              string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Test Accrual Route",
			method:           http.MethodGet,
			url:              "/api/orders/12345",
			expectedStatus:   http.StatusOK,
			expectedResponse: "accrual order",
		},
		{
			name:             "Test Register Order Route",
			method:           http.MethodPost,
			url:              "/api/orders",
			expectedStatus:   http.StatusOK,
			expectedResponse: "register order",
		},
		{
			name:             "Test Register Goods Route",
			method:           http.MethodPost,
			url:              "/api/goods",
			expectedStatus:   http.StatusOK,
			expectedResponse: "register goods",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLoggingMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			mockGzipMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			accrualHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("accrual order"))
			})
			registerOrderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("register order"))
			})
			registerGoodsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("register goods"))
			})

			router := NewAccrualRouter(
				mockLoggingMiddleware,
				mockGzipMiddleware,
				accrualHandler,
				registerOrderHandler,
				registerGoodsHandler,
			)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
