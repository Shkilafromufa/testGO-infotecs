package errors

import "errors"

// Общие ошибки, которые используются в приложении.
var (
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrSelfTransfer      = errors.New("you can't send yourself")
	ErrSenderNotFound    = errors.New("the sender's account does not exist")
	ErrRecipientNotFound = errors.New("the recipient's account does not exist")
	ErrInsufficientFunds = errors.New("the transfer amount is higher than the balance amount")
	ErrWalletNotFound    = errors.New("the wallet hash doesn't exist")
	ErrDbPathNotFound    = errors.New("the db path doesn't exist")
)
