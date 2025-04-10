package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		url              string
		authHeader       string
		expectedStatus   int
		expectedResponse string
		mockDecodeReturn *jwt.Claims
		mockDecodeError  error
	}{
		{
			name:             "Test Register Route",
			method:           http.MethodPost,
			url:              "/api/user/register",
			authHeader:       "",
			expectedStatus:   http.StatusOK,
			expectedResponse: "register",
			mockDecodeReturn: nil,
			mockDecodeError:  nil,
		},
		{
			name:             "Test Login Route",
			method:           http.MethodPost,
			url:              "/api/user/login",
			authHeader:       "",
			expectedStatus:   http.StatusOK,
			expectedResponse: "login",
			mockDecodeReturn: nil,
			mockDecodeError:  nil,
		},
		{
			name:             "Test Upload Order Route with Auth",
			method:           http.MethodPost,
			url:              "/api/user/orders",
			authHeader:       "Bearer some_token",
			expectedStatus:   http.StatusOK,
			expectedResponse: "upload order",
			mockDecodeReturn: &jwt.Claims{},
			mockDecodeError:  nil,
		},
		{
			name:             "Test Get Orders Route with Auth",
			method:           http.MethodGet,
			url:              "/api/user/orders",
			authHeader:       "Bearer some_token",
			expectedStatus:   http.StatusOK,
			expectedResponse: "get orders",
			mockDecodeReturn: &jwt.Claims{},
			mockDecodeError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockJWTDecoder := NewMockJWTDecoder(ctrl)
			mockLoggingMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			mockGzipMiddleware := func(next http.Handler) http.Handler {
				return next
			}
			mockAuthMiddleware := func(decoder JWTDecoder) func(next http.Handler) http.Handler {
				return func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						token := r.Header.Get("Authorization")
						_, err := decoder.Decode(token)
						if err != nil {
							http.Error(w, "Unauthorized", http.StatusUnauthorized)
							return
						}
						next.ServeHTTP(w, r)
					})
				}
			}

			registerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("register"))
			})
			loginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("login"))
			})
			uploadOrderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("upload order"))
			})
			getOrdersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("get orders"))
			})
			getBalanceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("balance"))
			})
			withdrawBalanceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("withdraw balance"))
			})
			withdrawalsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("withdrawals"))
			})

			router := NewGophermartRouter(
				mockLoggingMiddleware,
				mockGzipMiddleware,
				mockJWTDecoder,
				mockAuthMiddleware,
				registerHandler,
				loginHandler,
				uploadOrderHandler,
				getOrdersHandler,
				getBalanceHandler,
				withdrawBalanceHandler,
				withdrawalsHandler,
			)

			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			if tt.mockDecodeReturn != nil {
				mockJWTDecoder.EXPECT().Decode(tt.authHeader).Return(tt.mockDecodeReturn, tt.mockDecodeError).Times(1)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
