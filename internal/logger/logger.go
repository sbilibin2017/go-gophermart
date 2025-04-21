package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет собой глобальную переменную, содержащую настроенный экземпляр логгера.
// Логгер используется для записи логов с разными уровнями важности.
// Уровень логирования по умолчанию установлен на `InfoLevel`, что означает вывод сообщений
// с уровнем `Info` и выше.
var Logger *zap.Logger

// init инициализирует глобальный логгер `Logger` с использованием конфигурации по умолчанию.
func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	Logger, _ = zapConfig.Build()
}
