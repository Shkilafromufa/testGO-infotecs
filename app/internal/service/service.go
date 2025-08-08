package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	_ "net/rpc"
	"strconv"
	"testCaseGO/internal/model"
)

func Send(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newTrans := model.Transaction{}
	err = json.Unmarshal(body, &newTrans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if newTrans.From == newTrans.To {
		http.Error(w, "You can't send yourself", http.StatusBadRequest)
		return nil
	}
	_, err = ComparisonWallets(newTrans.From)
	if err != nil {
		http.Error(w, "The sender's account does not exist", http.StatusBadRequest)
		return nil
	}
	_, err = ComparisonWallets(newTrans.To)
	if err != nil {
		http.Error(w, "The recipient's account does not exist", http.StatusBadRequest)
		return nil
	}
	err = CompleteTransaction(newTrans.Amount, newTrans.From, newTrans.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	err = addNewTransaction(newTrans)
	if err != nil {
		http.Error(w, "The transaction add failed", http.StatusBadRequest)
		return err
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("The transaction was sent successfully"))
	if err != nil {
		return err
	}
	return nil
}
func GetLast(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	count := r.URL.Query().Get("count")
	countInt, err := strconv.Atoi(count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var transactions []model.Transaction
	transactions, err = GiveTransactionForCount(countInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	response, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	return nil
}
func GetBalance(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	walletHash := vars["address"]
	balance, err := GetBalanceFromDb(walletHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	response, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	return nil
}
