// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/clover-network/ethereum-spv-backend/app/controller"
	"github.com/gorilla/mux"
)

func verifyTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	txID := vars["id"]

	_, merklePath, err := controller.VerifyTransaction(txID)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]string)
	for k, v := range merklePath {
		value := fmt.Sprintf("0x%x", v)
		m["0x"+k] = value
	}

	json.NewEncoder(w).Encode(m)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/verify/{id}", verifyTransaction)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
