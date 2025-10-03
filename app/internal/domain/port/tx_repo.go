package port

import (
	"context"
	"testCaseGO/internal/domain/model"
)

// TxRepository описывает интерфейс работы с транзакциями.
type TxRepository interface {
	// Begin начинает новую БД-транзакцию и возвращает объект, зависящий от реализации.
	Begin(ctx context.Context) (any, error)
	// Commit завершает БД-транзакцию.
	Commit(ctx context.Context, tx any) error
	// Rollback откатывает БД-транзакцию.
	Rollback(ctx context.Context, tx any) error
	// Insert добавляет запись о новой транзакции в БД.
	Insert(ctx context.Context, tx any, t model.Transaction) error
	// GetLast возвращает последние N транзакций в порядке убывания даты.
	GetLast(ctx context.Context, limit int) ([]model.Transaction, error)
}
