package model

// Transaction описывает перевод между кошельками.
// Содержит адрес отправителя, получателя и сумму перевода.
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
