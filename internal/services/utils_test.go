package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToStruct(t *testing.T) {
	// Пример данных, которые мы будем маппить в структуру
	orderData := map[string]any{
		"order":   "123456",
		"status":  "PROCESSED",
		"accrual": 500,
	}

	// Инициализация структуры
	var orderResponse AccrualOrderGetResponse

	// Маппинг данных из карты в структуру
	err := mapToStruct(&orderResponse, orderData)

	// Проверка на отсутствие ошибки
	assert.Nil(t, err)

	// Проверка правильности значений в структуре
	assert.Equal(t, "123456", orderResponse.Order)
	assert.Equal(t, "PROCESSED", orderResponse.Status)
	assert.Equal(t, int64(500), orderResponse.Accrual)
}
