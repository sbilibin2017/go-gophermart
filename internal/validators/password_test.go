package usecases

import (
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator_Validate(t *testing.T) {
	// Создаем новый валидатор
	validator := NewPasswordValidator()

	// Массив тестов
	tests := []struct {
		password string
		expected error
	}{
		// Тесты с ошибками
		{"short", errors.ErrInvalidPassword},                       // слишком короткий пароль
		{"ThisIsAVeryLongPassword123!", errors.ErrInvalidPassword}, // слишком длинный пароль
		{"Password with space1!", errors.ErrInvalidPassword},       // пароль с пробелами
		{"password1!", errors.ErrInvalidPassword},                  // пароль без заглавной буквы
		{"Password!", errors.ErrInvalidPassword},                   // пароль без цифры
		{"Password1", errors.ErrInvalidPassword},                   // пароль без спецсимвола

		// Тест с успешной валидацией
		{"ValidPass1!", nil}, // правильный пароль
	}

	// Перебираем тесты
	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			// Проверяем результат валидации
			err := validator.Validate(tt.password)
			// Ожидаем, что ошибка совпадает с ожидаемой
			assert.Equal(t, tt.expected, err)
		})
	}
}
