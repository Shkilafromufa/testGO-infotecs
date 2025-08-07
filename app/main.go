package main

import (
	"net/http"
	"testCaseGO/internal/handler"
	"testCaseGO/internal/service"
)

func main() {
	handler.HandleRequest()
	service.CreateTables()
	service.ExistsWallets()
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		return
	}

}
