package sqlite

import (
	"context"
	"database/sql"
	appErr "testCaseGO/internal/errors"

	"testCaseGO/internal/domain/port"
)

// WalletRepo реализация интерфейса WalletRepository для SQLite.
type WalletRepo struct {
	db *DB
}

// NewWalletRepo создаёт новый WalletRepo.
func NewWalletRepo(db *DB) port.WalletRepository { return &WalletRepo{db: db} }

// GetBalance возвращает баланс кошелька по его адресу (hash).
func (r *WalletRepo) GetBalance(ctx context.Context, hash string) (float64, error) {
	var b float64
	err := r.db.QueryRowContext(ctx, "SELECT Balance FROM wallet WHERE Hash = ?", hash).Scan(&b)
	if err != nil {
		return 0, appErr.ErrWalletNotFound
	}
	return b, nil
}

// Exists проверяет существование кошелька по его адресу.
func (r *WalletRepo) Exists(ctx context.Context, hash string) (bool, error) {
	var h string
	err := r.db.QueryRowContext(ctx, "SELECT Hash FROM wallet WHERE Hash = ?", hash).Scan(&h)
	if err != nil {
		return false, err
	}
	return h == hash, nil
}

// Debit уменьшает баланс кошелька на указанную сумму в рамках SQL-транзакции.
func (r *WalletRepo) Debit(ctx context.Context, tx any, hash string, amount float64) error {
	sqltx := tx.(*sql.Tx)
	var bal float64
	if err := sqltx.QueryRowContext(ctx, "SELECT Balance FROM wallet WHERE Hash = ?", hash).Scan(&bal); err != nil {
		return err
	}
	if bal-amount < 0 {
		return appErr.ErrInsufficientFunds
	}
	_, err := sqltx.ExecContext(ctx, "UPDATE wallet SET Balance = ? WHERE Hash = ? AND Balance = ?", bal-amount, hash, bal)
	return err
}

// Credit увеличивает баланс кошелька на указанную сумму в рамках SQL-транзакции.
func (r *WalletRepo) Credit(ctx context.Context, tx any, hash string, amount float64) error {
	sqltx := tx.(*sql.Tx)
	var bal float64
	if err := sqltx.QueryRowContext(ctx, "SELECT Balance FROM wallet WHERE Hash = ?", hash).Scan(&bal); err != nil {
		return err
	}
	_, err := sqltx.ExecContext(ctx, "UPDATE wallet SET Balance = ? WHERE Hash = ? AND Balance = ?", bal+amount, hash, bal)
	return err
}
