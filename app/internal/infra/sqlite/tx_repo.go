package sqlite

import (
	"context"
	"database/sql"
	"testCaseGO/internal/domain/model"
	"testCaseGO/internal/domain/port"
)

// TxRepo реализация интерфейса TxRepository для SQLite.
type TxRepo struct {
	db *DB
}

// NewTxRepo создаёт новый TxRepo.
func NewTxRepo(db *DB) port.TxRepository { return &TxRepo{db: db} }

// Begin начинает SQL-транзакцию и возвращает её объект.
func (r *TxRepo) Begin(ctx context.Context) (any, error) {
	return r.db.BeginTx(ctx, nil)
}

// Commit выполняет commit SQL-транзакции.
func (r *TxRepo) Commit(ctx context.Context, tx any) error {
	return tx.(*sql.Tx).Commit()
}

// Rollback выполняет rollback SQL-транзакции.
func (r *TxRepo) Rollback(ctx context.Context, tx any) error {
	return tx.(*sql.Tx).Rollback()
}

// Insert добавляет новую запись в таблицу transactions в рамках транзакции.
func (r *TxRepo) Insert(ctx context.Context, tx any, t model.Transaction) error {
	_, err := tx.(*sql.Tx).ExecContext(ctx,
		"INSERT INTO transactions(From_address, To_address, Amount) VALUES (?,?,?)",
		t.From, t.To, t.Amount,
	)
	return err
}

// GetLast возвращает последние limit транзакций, отсортированные по Timestamp DESC.
func (r *TxRepo) GetLast(ctx context.Context, limit int) ([]model.Transaction, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT From_address, To_address, Amount FROM transactions ORDER BY Timestamp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.From, &t.To, &t.Amount); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}
