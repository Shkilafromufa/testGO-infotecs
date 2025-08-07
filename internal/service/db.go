package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testCaseGO/internal/model"
)

func CreateTables() {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)
	// Создаем первую таблицу
	prepare, err := db.Prepare("CREATE TABLE IF NOT EXISTS transactions (Id INTEGER PRIMARY KEY AUTOINCREMENT, Amount INTEGER, From_address TEXT, To_address TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = prepare.Exec()
	if err != nil {
		log.Fatal(err)
	}
	// Создаем вторую таблицу
	prepare, err = db.Prepare("CREATE TABLE IF NOT EXISTS wallet(Id INTEGER PRIMARY KEY AUTOINCREMENT, Hash TEXT, Balance INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = prepare.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
func initDb() (*sql.DB, error) {
	bd, err := sql.Open("sqlite3", "../internal/service/database/database.db")
	if err != nil {
		return nil, err
	}
	return bd, err
}
func InitStartWallets() {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)
	prepare, err := db.Prepare("INSERT INTO wallet(Hash, Balance) VALUES (?, ?)")
	if err != nil {
		return
	}
	var count = 10
	for i := 0; i < count; i++ {
		var buf [32]byte
		var wallet model.Wallet
		if _, err := rand.Read(buf[:]); err != nil {
			log.Fatalf("rand.Read: %v", err)
		}
		wallet.Hash = hex.EncodeToString(buf[:])
		wallet.Balance = 100
		_, err := prepare.Exec(wallet.Hash, wallet.Balance)
		if err != nil {
			return
		}
	}
}
func ExistsWallets() {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db)
	var icon = 0 // переменная "маячок"
	// Проверяем создавались ли уже ранее кошельки
	err = db.QueryRow("SELECT Id FROM wallet WHERE Id = ?", 1).Scan(&icon)
	if icon != 1 {
		InitStartWallets()
	}
}
func addNewTransaction(transaction model.Transaction) error {
	db, err := initDb()
	if err != nil {
		return err
	}
	defer closeDB(db)
	prepare, err := db.Prepare("INSERT INTO transactions(From_address, To_address, Amount) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = prepare.Exec(transaction.From, transaction.To, transaction.Amount)
	return nil
}
func ComparisonWallets(address string) (bool, error) {
	db, err := initDb()
	if err != nil {
		return false, err
	}
	defer closeDB(db)
	var hashes string
	err = db.QueryRow("SELECT Hash FROM wallet WHERE Hash = ?", address).Scan(&hashes)
	if err != nil {
		return false, err
	}
	if hashes == address {
		return true, nil
	}
	return false, fmt.Errorf("address %s not found in wallets", address)
}
func CompleteTransaction(amount float64, hashFrom string, hashTo string) error {
	db, err := initDb()
	if err != nil {
		return err
	}
	var BalanceFrom, BalanceTo float64
	defer closeDB(db)
	err = db.QueryRow("SELECT Balance FROM wallet WHERE Hash =?", hashFrom).Scan(&BalanceFrom)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT Balance FROM wallet WHERE Hash =?", hashTo).Scan(&BalanceTo)
	if err != nil {
		return err
	}
	var newBalanceFrom float64
	var newBalanceTo float64
	if BalanceFrom-amount < 0 {
		return errors.New(`the transfer amount is higher than the balance amount`)
	}
	newBalanceFrom = BalanceFrom - amount
	newBalanceTo = BalanceTo + amount
	prepare, err := db.Prepare("UPDATE wallet SET Balance=? WHERE Hash=? AND Balance=?")
	if err != nil {
		return err
	}
	_, err = prepare.Exec(newBalanceFrom, hashFrom, BalanceFrom)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(newBalanceTo, hashTo, BalanceTo)
	if err != nil {
		return err
	}
	return nil
}
func GiveTransactionForCount(count int) ([]model.Transaction, error) {
	db, err := initDb()
	if err != nil {
		return nil, err
	}
	defer closeDB(db)
	state, err := db.Prepare("SELECT From_address, To_address, Amount FROM transactions ORDER BY Timestamp DESC LIMIT ?")
	if err != nil {
		return nil, err
	}
	rows, err := state.Query(count)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(&t.From, &t.To, &t.Amount)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, err
}
func GetBalanceFromDb(hash string) (float64, error) {
	db, err := initDb()
	if err != nil {
		return 0, err
	}
	defer closeDB(db)
	var balanceFrom float64
	state, err := db.Prepare("SELECT Balance FROM wallet WHERE Hash=?")
	if err != nil {
		return 0, err
	}
	rows := state.QueryRow(hash)
	err = rows.Scan(&balanceFrom)
	if err != nil {
		err = errors.New(`the wallet hash doesn't exist'`)
		return 0, err
	}
	return balanceFrom, nil
}
func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}
