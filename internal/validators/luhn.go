package validators

import (
	"strconv"
	"strings"
)

// LuhnCheck выполняет проверку числа по алгоритму Луна
func ValidateNumberWithLuhn(number string) bool {
	// Убираем все пробелы из строки
	number = strings.ReplaceAll(number, " ", "")

	// Проверяем, что строка состоит только из цифр
	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}

	// Преобразуем строку в массив чисел
	sum := 0
	shouldDouble := false

	// Идем с конца по числу
	for i := len(number) - 1; i >= 0; i-- {
		// Преобразуем текущий символ в число
		n, _ := strconv.Atoi(string(number[i]))

		// Удваиваем каждую вторую цифру с конца
		if shouldDouble {
			n *= 2
			// Если результат больше 9, то складываем его цифры
			if n > 9 {
				n -= 9
			}
		}

		// Добавляем значение к сумме
		sum += n

		// Переключаем флаг для удвоения цифры
		shouldDouble = !shouldDouble
	}

	// Проверяем, делится ли сумма на 10
	return sum%10 == 0
}
