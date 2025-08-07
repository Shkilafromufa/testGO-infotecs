package handler

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testCaseGO/internal/service"
)

func HandleRequest() {
	http.HandleFunc("/api/transactions", HandleLastTrans)
	http.HandleFunc("/api/send", HandleSendCash)
	balanceRouter := mux.NewRouter()
	balanceRouter.Methods(http.MethodGet).Path("/api/wallet/{address}/balance").HandlerFunc(HandleBalance)
	http.Handle("/api/wallet/", balanceRouter)
}

func HandleLastTrans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			err := service.GetLast(w, r)
			if err != nil {
				log.Println(err)
			}
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}
func HandleBalance(w http.ResponseWriter, r *http.Request) {
	err := service.GetBalance(w, r)
	if err != nil {
		log.Println(err)
	}
}
func HandleSendCash(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			err := service.Send(w, r)
			if err != nil {
				log.Println(err)
			}
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
