package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"testCaseGO/internal/domain/model"
	"testCaseGO/internal/service"
)

// HTTP отвечает за HTTP-слой (хендлеры Fiber).
type HTTP struct {
	svc *service.TransferService
}

// NewHTTP регистрирует HTTP-роуты в Fiber и привязывает их к хендлерам.
func NewHTTP(app *fiber.App, svc *service.TransferService) {
	h := &HTTP{svc: svc}
	api := app.Group("/api")
	api.Get("/transactions", h.getLast)
	api.Post("/send", h.send)
	api.Get("/wallet/:address/balance", h.balance)
}

// send обрабатывает POST /api/send.
// Принимает JSON с переводом, вызывает бизнес-логику и возвращает результат.
func (h *HTTP) send(c *fiber.Ctx) error {
	var t model.Transaction
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.svc.Send(c.Context(), t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "The transaction was sent successfully"})
}

// getLast обрабатывает GET /api/transactions.
// Возвращает список последних транзакций, количество задаётся через query-параметр count.
func (h *HTTP) getLast(c *fiber.Ctx) error {
	limitStr := c.Query("count", "10")
	limit, _ := strconv.Atoi(limitStr)
	list, err := h.svc.Last(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(list)
}

// balance обрабатывает GET /api/wallet/:address/balance.
// Возвращает баланс кошелька по его адресу.
func (h *HTTP) balance(c *fiber.Ctx) error {
	hash := c.Params("address")
	bal, err := h.svc.Balance(c.Context(), hash)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(bal)
}
