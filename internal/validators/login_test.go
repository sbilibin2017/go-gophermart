package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValidator_Validate(t *testing.T) {
	// Создаем экземпляр валидатора
	validator := NewLoginValidator()

	// Сценарии с корректными логинами
	validLogins := []string{
		"user123",
		"valid_user_1",
		"username-123",
	}

	for _, login := range validLogins {
		t.Run("valid login: "+login, func(t *testing.T) {
			// Проверяем, что для каждого валидного логина не возникает ошибки
			err := validator.Validate(login)
			assert.NoError(t, err, "Expected no error for valid login")
		})
	}

	// Сценарии с некорректными логинами
	invalidLogins := []struct {
		login       string
		expectedErr error
	}{
		{"ab", ErrInvalidLogin},                       // Логин слишком короткий
		{"thisisaverylongloginname", ErrInvalidLogin}, // Логин слишком длинный
		{"user name", ErrInvalidLogin},                // Логин содержит пробел
		{"user@name", ErrInvalidLogin},                // Логин содержит недопустимый символ '@'
		{"user#123", ErrInvalidLogin},                 // Логин содержит недопустимый символ '#'
		{"тест", ErrInvalidLogin},
	}

	for _, test := range invalidLogins {
		t.Run("invalid login: "+test.login, func(t *testing.T) {
			// Проверяем, что для каждого некорректного логина возникает соответствующая ошибка
			err := validator.Validate(test.login)
			assert.NotNil(t, err)
		})
	}
}
