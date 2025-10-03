package port

import "context"

// WalletRepository описывает интерфейс работы с кошельками.
type WalletRepository interface {
	// GetBalance возвращает баланс по хэшу кошелька.
	GetBalance(ctx context.Context, hash string) (float64, error)
	// Exists проверяет, существует ли кошелёк с данным хэшем.
	Exists(ctx context.Context, hash string) (bool, error)
	// Debit уменьшает баланс кошелька на указанную сумму в рамках транзакции.
	Debit(ctx context.Context, tx any, hash string, amount float64) error
	// Credit увеличивает баланс кошелька на указанную сумму в рамках транзакции.
	Credit(ctx context.Context, tx any, hash string, amount float64) error
}
