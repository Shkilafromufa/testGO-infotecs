package setup

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"testCaseGO/internal/infra/sqlite"
	"testCaseGO/internal/service"
)

// Container содержит зависимости приложения.
type Container struct {
	DB  *sqlite.DB
	Svc *service.TransferService
}

// PrepareEnv инициализирует окружение: загружает .env,
// открывает БД, выполняет миграции, добавляет данные и возвращает контейнер зависимостей.
func PrepareEnv() (*Container, error) {
	_ = godotenv.Load("app/internal/config/.env")
	if os.Getenv("DBPATH") == "" {
		log.Println("internal/service/database/database.db")
	}

	db, err := sqlite.NewDB()
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(); err != nil {
		db.Close()
		return nil, err
	}
	if err := db.SeedWalletsIfEmpty(); err != nil {
		db.Close()
		return nil, err
	}

	wr := sqlite.NewWalletRepo(db)
	tr := sqlite.NewTxRepo(db)
	svc := service.NewTransferService(wr, tr)

	log.Println("Environment ready")
	return &Container{DB: db, Svc: svc}, nil
}
