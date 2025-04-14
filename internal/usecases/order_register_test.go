package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestOrderRegisterUsecase_Execute(t *testing.T) {
	tests := []struct {
		name                  string
		req                   *types.OrderRegisterRequest
		validateErr           error
		registerErr           error
		expectedResp          *types.OrderRegisterResponse
		expectedErr           error
		expectedRegisterTimes int
	}{
		{
			name: "Valid Request",
			req: &types.OrderRegisterRequest{
				Number: 123,
				Goods: []types.Good{
					{Description: "Item 1", Price: 100},
					{Description: "Item 2", Price: 200},
				},
			},
			validateErr: nil,
			registerErr: nil,
			expectedResp: &types.OrderRegisterResponse{
				Message: "order registered successfully",
			},
			expectedErr:           nil,
			expectedRegisterTimes: 1,
		},
		{
			name: "Invalid Request - Validation Error",
			req: &types.OrderRegisterRequest{
				Number: 123,
				Goods: []types.Good{
					{Description: "Item 1", Price: 100},
					{Description: "", Price: 200},
				},
			},
			validateErr:           ErrOrderRegisterInvalidRequest,
			registerErr:           nil,
			expectedResp:          nil,
			expectedErr:           ErrOrderRegisterInvalidRequest,
			expectedRegisterTimes: 0,
		},
		{
			name: "Service Failure",
			req: &types.OrderRegisterRequest{
				Number: 123,
				Goods: []types.Good{
					{Description: "Item 1", Price: 100},
					{Description: "Item 2", Price: 200},
				},
			},
			validateErr:           nil,
			registerErr:           errors.New("service error"),
			expectedResp:          nil,
			expectedErr:           errors.New("service error"),
			expectedRegisterTimes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockValidator := NewMockOrderValidator(ctrl)
			mockService := NewMockOrderRegisterService(ctrl)

			uc := NewOrderRegisterUsecase(mockValidator, mockService)

			mockValidator.EXPECT().Struct(tt.req).Return(tt.validateErr).Times(1)
			mockService.EXPECT().Register(gomock.Any(), gomock.Eq(&types.Order{
				Number: tt.req.Number,
				Goods: []types.Good{
					{Description: "Item 1", Price: 100},
					{Description: "Item 2", Price: 200},
				},
			})).Return(tt.registerErr).Times(tt.expectedRegisterTimes)

			resp, err := uc.Execute(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
