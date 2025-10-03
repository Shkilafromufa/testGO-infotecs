package model

// Wallet представляет кошелёк в системе.
// Hash — уникальный идентификатор (адрес).
// Balance — текущий баланс кошелька.
type Wallet struct {
	Hash    string `json:"hash"`
	Balance int64  `json:"balance"`
}
