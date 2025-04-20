package repositories

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestExecContext_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем успешное выполнение запроса
	mock.ExpectExec("UPDATE rewards").
		WithArgs("value").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создаем функцию txProvider
	txProvider := func(ctx context.Context) *sql.Tx {
		return nil // Нет транзакции в данном случае
	}

	// Вызов метода execContext
	result, err := execContext(context.Background(), db, txProvider, "UPDATE rewards SET field = ?", "value")

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что количество измененных строк равно 1
	rowsAffected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), rowsAffected)
}

func TestExecContext_Error(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем ошибку выполнения запроса
	mock.ExpectExec("UPDATE rewards").
		WithArgs("value").
		WillReturnError(sql.ErrConnDone)

	// Создаем функцию txProvider
	txProvider := func(ctx context.Context) *sql.Tx {
		return nil // Нет транзакции в данном случае
	}

	// Вызов метода execContext
	result, err := execContext(context.Background(), db, txProvider, "UPDATE rewards SET field = ?", "value")

	// Проверяем, что ошибка возникла
	assert.Error(t, err)

	// Проверяем, что результат nil, так как произошла ошибка
	assert.Nil(t, result)
}

func TestQueryContext_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем успешное выполнение запроса
	mock.ExpectQuery("SELECT 1").
		WithArgs("match1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(1))

	// Создаем функцию txProvider
	txProvider := func(ctx context.Context) *sql.Tx {
		return nil // Нет транзакции в данном случае
	}

	// Вызов метода queryContext
	row, err := queryRowContext(context.Background(), db, txProvider, "SELECT EXISTS(SELECT 1 FROM rewards WHERE match = ?)", "match1")

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что row не nil
	assert.NotNil(t, row)

	var exists int
	err = row.Scan(&exists)
	assert.NoError(t, err)
	assert.Equal(t, 1, exists)
}

func TestQueryContext_Error(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем ошибку выполнения запроса
	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM rewards WHERE match = ?)").
		WithArgs("match1").
		WillReturnError(sql.ErrConnDone) // Симулируем ошибку

	// Создаем функцию txProvider
	txProvider := func(ctx context.Context) *sql.Tx {
		return nil // Нет транзакции в данном случае
	}

	// Вызов метода queryContext
	row, err := queryRowContext(context.Background(), db, txProvider, "SELECT EXISTS(SELECT 1 FROM rewards WHERE match = ?)", "match1")

	// Проверяем, что ошибка возникла
	assert.Error(t, err)

	// Проверяем, что row == nil, так как произошла ошибка
	assert.Nil(t, row)
}

func TestQueryContext_WithTx(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем успешное выполнение запроса с транзакцией
	mock.ExpectBegin()
	mock.ExpectQuery("").
		WithArgs("match1").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectCommit()

	// Мокаем транзакцию
	txProvider := func(ctx context.Context) *sql.Tx {
		tx, _ := db.Begin()
		return tx
	}

	// Вызов метода queryContext с транзакцией
	row, err := queryRowContext(context.Background(), db, txProvider, "SELECT EXISTS(SELECT 1 FROM rewards WHERE match = ?)", "match1")

	// Проверяем, что ошибок нет
	assert.NoError(t, err)
	// Проверяем, что row != nil, так как запрос должен быть успешным
	assert.NotNil(t, row)
}

func TestExecContext_WithTx(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания mock БД: %v", err)
	}
	defer db.Close()

	// Мокаем успешное выполнение запроса с транзакцией
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE rewards SET match = ?").
		WithArgs("new_match1").
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 - id строки, 1 - количество затронутых строк
	mock.ExpectCommit()

	// Мокаем транзакцию
	txProvider := func(ctx context.Context) *sql.Tx {
		tx, _ := db.Begin()
		return tx
	}

	// Вызов метода execContext с транзакцией
	result, err := execContext(context.Background(), db, txProvider, "UPDATE rewards SET match = ?", "new_match1")

	// Проверяем, что ошибок нет
	assert.NoError(t, err)
	// Проверяем, что result != nil, так как запрос должен быть успешным
	assert.NotNil(t, result)
}
