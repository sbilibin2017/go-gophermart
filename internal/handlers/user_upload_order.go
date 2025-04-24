package handlers

// import (
// 	"context"
// 	"net/http"

// 	"github.com/sbilibin2017/go-gophermart/internal/api/middlewares"
// 	"github.com/sbilibin2017/go-gophermart/internal/services/domain"
// )

// type UploadUserOrderService interface {
// 	Upload(ctx context.Context, userOrder *domain.UserOrder) error
// }

// type UploadUserOrderValidator interface {
// 	Struct(v any) error
// }

// type UploadUserOrderRequest struct {
// 	Number string `json:"number" validate:"required,luhn"`
// }

// func UploadUserOrderHandler(svc UploadUserOrderService, v UploadUserOrderValidator) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		login, err := middlewares.GetLogin(r)
// 		if err != nil {
// 			handleUnauthorizedError(w, err)
// 			return
// 		}

// 		number := getPathParam(r, "number")

// 		err = validateRequest(UploadUserOrderRequest{Number: number}, v)
// 		if err != nil {
// 			handleBadRequestError(w, err)
// 			return
// 		}

// 		err = svc.Upload(r.Context(), &domain.UserOrder{UserID: login, Order: domain.Order{OrderID: number}})

// 		if err != nil {
// 			switch err {
// 			case domain.ErrOrderAlreadyUploadedByUser:
// 				writeTextPlainResponse(w, "Order number already uploaded by this user", http.StatusOK)
// 				return
// 			case domain.ErrOrderAlreadyUploadedByAnotherUser:
// 				handleConflictError(w, err)
// 				return
// 			default:
// 				handleInternalError(w, err)
// 				return
// 			}
// 		}

// 		writeTextPlainResponse(w, "Order number successfully uploaded", http.StatusAccepted)
// 	}
// }
