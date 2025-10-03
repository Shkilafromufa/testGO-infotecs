package service

import (
	"context"
	"testCaseGO/internal/domain/model"
	"testCaseGO/internal/domain/port"
	appErr "testCaseGO/internal/errors"
)

// TransferService реализует бизнес-логику переводов и работы с кошельками.
type TransferService struct {
	wallet port.WalletRepository
	txRepo port.TxRepository
}

// NewTransferService создаёт новый сервис переводов.
func NewTransferService(w port.WalletRepository, t port.TxRepository) *TransferService {
	return &TransferService{wallet: w, txRepo: t}
}

// Send выполняет перевод между кошельками.
// Сначала проверяются все условия (сумма > 0, разные адреса, наличие обоих кошельков).
// Затем выполняется транзакция: списание, начисление и запись истории.
func (s *TransferService) Send(ctx context.Context, t model.Transaction) error {
	if t.From == t.To {
		return appErr.ErrSelfTransfer
	}
	if t.Amount <= 0 {
		return appErr.ErrInvalidAmount
	}

	ok, err := s.wallet.Exists(ctx, t.From)
	if err != nil || !ok {
		return appErr.ErrSenderNotFound
	}
	ok, err = s.wallet.Exists(ctx, t.To)
	if err != nil || !ok {
		return appErr.ErrRecipientNotFound
	}

	tx, err := s.txRepo.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(txRepo port.TxRepository, ctx context.Context, tx any) {
		err := txRepo.Rollback(ctx, tx)
		if err != nil {

		}
	}(s.txRepo, ctx, tx)

	if err := s.wallet.Debit(ctx, tx, t.From, t.Amount); err != nil {
		return err
	}
	if err := s.wallet.Credit(ctx, tx, t.To, t.Amount); err != nil {
		return err
	}
	if err := s.txRepo.Insert(ctx, tx, t); err != nil {
		return err
	}

	return s.txRepo.Commit(ctx, tx)
}

// Last возвращает последние N транзакций (по умолчанию 10).
func (s *TransferService) Last(ctx context.Context, limit int) ([]model.Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.txRepo.GetLast(ctx, limit)
}

// Balance возвращает баланс по хэшу кошелька.
func (s *TransferService) Balance(ctx context.Context, hash string) (float64, error) {
	return s.wallet.GetBalance(ctx, hash)
}
