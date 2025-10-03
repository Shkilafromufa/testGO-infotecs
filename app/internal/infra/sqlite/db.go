package sqlite

import (
	"database/sql"
	"os"
	appErr "testCaseGO/internal/errors"

	_ "github.com/mattn/go-sqlite3"
)

// DB обёртка над *sql.DB, используемая в приложении.
// Содержит соединение с SQLite и методы для миграций/сидирования.
type DB struct{ *sql.DB }

// NewDB открывает соединение с SQLite, используя путь из переменной окружения DBPATH.
// Если DBPATH не задана, используется дефолтный путь internal/service/database/database.db.
// Возвращает обёртку DB и ошибку.
func NewDB() (*DB, error) {
	path := os.Getenv("DBPATH")
	if path == "" {
		return nil, appErr.ErrDbPathNotFound
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Migrate выполняет миграции схемы БД: создаёт таблицы transactions и wallet, если их нет.
func (d *DB) Migrate() error {
	_, err := d.Exec(`CREATE TABLE IF NOT EXISTS transactions (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Amount INTEGER,
		From_address TEXT,
		To_address TEXT,
		Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	_, err = d.Exec(`CREATE TABLE IF NOT EXISTS wallet (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Hash TEXT UNIQUE,
		Balance INTEGER
	)`)
	return err
}

// SeedWalletsIfEmpty заполняет таблицу wallet начальными данными, если она пуста.
// Создаются 10 кошельков с начальными балансами по 100.
func (d *DB) SeedWalletsIfEmpty() error {
	var id int
	err := d.QueryRow("SELECT Id FROM wallet LIMIT 1").Scan(&id)
	if err == nil { // уже есть
		return nil
	}
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO wallet(Hash, Balance) VALUES (?,?)")
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	for i := 0; i < 10; i++ {
		hash := fmtHash(i)
		if _, err := stmt.Exec(hash, 100); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func fmtHash(i int) string {
	return "w" + fmtInt(i)
}
func fmtInt(i int) string {
	if i < 10 {
		return "0" + string('0'+byte(i))
	}
	return string('0'+byte(i/10)) + string('0'+byte(i%10))
}

func (d *DB) Close() { _ = d.DB.Close() }
