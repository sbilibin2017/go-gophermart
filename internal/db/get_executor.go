package db

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

// GetExecutor извлекает или создаёт исполнитель для работы с базой данных в зависимости от контекста.
//
// Эта функция проверяет, содержится ли транзакция в контексте (через функцию GetTx). Если транзакция найдена,
// то возвращается экземпляр типа TxExecutor, который работает с транзакцией. Если транзакция не найдена,
// возвращается экземпляр типа DBExecutor, который работает с обычным соединением с базой данных.
//
// Параметры:
//   - ctx: контекст, в котором может быть сохранена транзакция.
//   - db: объект базы данных (*sql.DB), который будет использован, если транзакция не найдена в контексте.
//
// Возвращает:
//   - Executor: интерфейс, который позволяет выполнять запросы и операции как в рамках транзакции, так и в обычной базе данных.
//     Возвращаемый объект может быть либо `*TxExecutor`, если в контексте есть транзакция, либо `*DBExecutor`, если транзакция отсутствует.
func GetExecutor(ctx context.Context, db *sql.DB) Executor {
	if tx, ok := GetTx(ctx); ok {
		return &txExecutor{tx: tx}
	}
	return &dbExecutor{db: db}
}

// Executor - интерфейс, который определяет базовые методы для работы с базой данных в контексте SQL-запросов.
// Он предоставляет методы для выполнения запросов и команд на уровне базы данных и транзакций.
//
// Интерфейс предоставляет абстракцию для выполнения запросов как в контексте обычной базы данных (DB),
// так и в контексте транзакции (TX). Методы интерфейса логируют выполняемые SQL-запросы,
// что позволяет отслеживать и анализировать взаимодействие с базой данных.

type Executor interface {
	// QueryRowContext выполняет SQL-запрос, который ожидает вернуть одну строку данных.
	// Возвращает *sql.Row, который позволяет получить результат из запроса.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	// QueryContext выполняет SQL-запрос и возвращает множество строк.
	// Возвращает *sql.Rows для чтения результата запроса.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// ExecContext выполняет SQL-команду, которая не возвращает строк (например, UPDATE, INSERT).
	// Возвращает sql.Result, который позволяет узнать, сколько строк было затронуто.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// dbExecutor - реализация интерфейса Executor для работы с обычной базой данных.
// Эта структура используется для выполнения запросов в контексте обычной базы данных, без использования транзакций.
// Каждый запрос логируется с помощью функции logQuery, которая записывает подробности запроса.
//
// Примечание: Логирование происходит до выполнения запроса, что позволяет отслеживать запросы
// в момент их вызова в коде.

type dbExecutor struct {
	db *sql.DB // Экземпляр базы данных для выполнения запросов
}

// logQuery записывает информацию о выполняемом SQL-запросе в лог.
func (e *dbExecutor) logQuery(msg, query string, args ...any) {
	// Логирование запроса с деталями
	logger.Logger.Info(msg, zap.String("query", query), zap.Any("args", args))
}

// QueryRowContext выполняет SQL-запрос, ожидая одну строку в результате.
func (e *dbExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.logQuery("Executing QueryRowContext", query, args...)
	return e.db.QueryRowContext(ctx, query, args...)
}

// QueryContext выполняет SQL-запрос, который возвращает несколько строк.
func (e *dbExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.logQuery("Executing QueryContext", query, args...)
	return e.db.QueryContext(ctx, query, args...)
}

// ExecContext выполняет SQL-команду, которая не возвращает строк.
func (e *dbExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.logQuery("Executing ExecContext", query, args...)
	return e.db.ExecContext(ctx, query, args...)
}

// txExecutor - реализация интерфейса Executor для работы с транзакциями.
// Эта структура используется для выполнения запросов в контексте транзакций.
// Также логируются все SQL-запросы, которые выполняются в транзакционном контексте.
//
// Примечание: Логирование выполняется до выполнения запроса, что позволяет отслеживать запросы,
// когда они выполняются в рамках транзакции.

type txExecutor struct {
	tx *sql.Tx // Экземпляр транзакции для выполнения запросов
}

// logQuery записывает информацию о выполняемом SQL-запросе в лог.
func (e *txExecutor) logQuery(msg, query string, args ...any) {
	// Логирование запроса с деталями
	logger.Logger.Info(msg, zap.String("query", query), zap.Any("args", args))
}

// QueryRowContext выполняет SQL-запрос, ожидая одну строку в результате.
func (e *txExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.logQuery("Executing QueryRowContext (TX)", query, args...)
	return e.tx.QueryRowContext(ctx, query, args...)
}

// QueryContext выполняет SQL-запрос, который возвращает несколько строк.
func (e *txExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.logQuery("Executing QueryContext (TX)", query, args...)
	return e.tx.QueryContext(ctx, query, args...)
}

// ExecContext выполняет SQL-команду, которая не возвращает строк.
func (e *txExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.logQuery("Executing ExecContext (TX)", query, args...)
	return e.tx.ExecContext(ctx, query, args...)
}
