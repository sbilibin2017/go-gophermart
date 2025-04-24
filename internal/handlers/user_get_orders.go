package handlers

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"github.com/sbilibin2017/go-gophermart/internal/api/middlewares"
// 	"github.com/sbilibin2017/go-gophermart/internal/services/domain"
// )

// type UserOrderResponse struct {
// 	Number     string `json:"number"`
// 	Status     string `json:"status"`
// 	Accrual    *int64 `json:"accrual,omitempty"`
// 	UploadedAt string `json:"uploaded_at"`
// }

// type UserGetOrdersService interface {
// 	Get(ctx context.Context, userID string) ([]*domain.UserOrder, error)
// }

// func GetUserOrdersHandler(svc UserGetOrdersService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// 401 — пользователь не авторизован
// 		login, err := middlewares.GetLogin(r)
// 		if err != nil {
// 			handleUnauthorizedError(w, err)
// 			return
// 		}

// 		orders, err := svc.Get(r.Context(), login)
// 		if err != nil {
// 			// 500 — внутренняя ошибка сервера
// 			handleInternalError(w, err)
// 			return
// 		}

// 		if len(orders) == 0 {
// 			// 204 — нет данных для ответа
// 			writeTextPlainResponse(w, "No orders found", http.StatusNoContent)
// 			return
// 		}

// 		var orderResponses []UserOrderResponse
// 		for _, order := range orders {
// 			orderResponses = append(orderResponses, UserOrderResponse{
// 				Number:     order.Number,
// 				Status:     string(order.Status),
// 				Accrual:    order.Accrual,
// 				UploadedAt: order.UploadedAt.Format(time.RFC3339),
// 			})
// 		}

// 		// 200 — успешная обработка запроса
// 		encodeJSONResponse(w, orderResponses, http.StatusOK)
// 	}
// }
